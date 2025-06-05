// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CascadingBloomFilter {
    /* ─── State ───────────────────────────────────────────────────────────── */

    address private immutable _owner;

    struct Layer {
        uint64  filterSizeBits;
        uint32  k;
        bytes   filter;
    }

    Layer[] private layers;

    /* ─── Modifiers ───────────────────────────────────────────────────────── */

    modifier onlyOwner() {
        require(msg.sender == _owner, "Not owner");
        _;
    }

    /* ─── Constructor ─────────────────────────────────────────────────────── */

    constructor() {
        _owner = msg.sender;
    }

    /* ─── Public/External API ──────────────────────────────────────────────── */

    /// @notice Replace all layers in one call.
    /// @param newFilters   newFilters[i] is the full packed bit‐vector for layer i
    /// @param ks           ks[i] = number of hash functions for layer i
    /// @param bitLens      bitLens[i] = number of bits used in layer i
    function updateCascade(
        bytes[] calldata newFilters,
        uint256[] calldata ks,
        uint256[] calldata bitLens
    ) external onlyOwner {
        uint256 len = newFilters.length;
        require(len > 0,             "At least one layer");
        require(len == ks.length,    "Need k for each layer");
        require(len == bitLens.length, "Need bitLen for each layer");

        // Wipe out existing layers (cheapest way to reset a dynamic array)
        delete layers;

        // Build up each Layer exactly once:
        for (uint256 i = 0; i < len; ) {
            bytes calldata f   = newFilters[i];
            uint256      k_   = ks[i];
            uint256      bits = bitLens[i];

            // -- validate inputs --
            require(k_ > 0,                     "k must be > 0");
            require(bits > 0,                   "bitLen must be > 0");
            require(bits <= f.length * 8,      "bitLen exceeds f.length*8");
            require(k_ <= type(uint32).max,    "k too large for uint32");
            require(bits <= type(uint64).max,  "bitLen too large for uint64");

            // Pack into (uint64 filterSizeBits, uint32 k, bytes filter)
            layers.push(
                Layer({
                    filterSizeBits: uint64(bits),
                    k:              uint32(k_),
                    filter:         f
                })
            );

            unchecked { ++i; }
        }
    }

    /// @notice Return the total number of layers.
    function layerCount() external view returns (uint256) {
        return layers.length;
    }

    /// @notice Test `token` against every layer. Returns (accepted, layerIndexReached).
    /// ‣ “Early‐accept” if you hit a zero‐bit in an odd‐indexed layer.
    /// ‣ On the last layer, require match == (lastIndex % 2 == 0).
    function testToken(bytes calldata token) external view returns (bool, uint256) {
        uint256 n = layers.length;
        require(n > 0, "No layers");

        // Precompute 4×64‐bit hashes once:
        uint64[4] memory h = extractHashes(token);

        for (uint256 li = 0; li < n; ) {
            Layer storage L = layers[li];
            bool match_     = _testInLayer(L.filter, uint256(L.filterSizeBits), L.k, h);

            // If this is the last layer:
            if (li == n - 1) {
                // On the final layer, we expect a “match” if and only if (li % 2 == 0).
                bool wantMatch = (li & 1) == 0;
                return (match_ == wantMatch, li);
            }

            if (!match_) {
                bool acceptEarly = (li & 1) == 1;
                return (acceptEarly, li);
            }

            unchecked { ++li; }
        }

        // This point should never happen.
        revert("unreachable");
    }

    /// @notice Return metadata and full filter bytes for layer i.
    function getLayerMetadata(uint256 i)
    external
    view
    returns (
        uint256 filterSizeBits_,
        uint256 k_,
        bytes memory filter_
    )
    {
        require(i < layers.length, "Invalid layer");
        Layer storage L = layers[i];
        return (uint256(L.filterSizeBits), uint256(L.k), L.filter);
    }

    /* ─── Internal Helpers ─────────────────────────────────────────────────── */

    /// @notice Test a single layer’s `filter` as a Bloom filter.
    /// @param filter         the packed bit‐vector (in storage)
    /// @param filterSizeBits number of bits of `filter` to consider
    /// @param k              number of hash probes
    /// @param h              4×64‐bit “precomputed” hashes
    function _testInLayer(
        bytes storage filter,
        uint256 filterSizeBits,
        uint256 k,
        uint64[4] memory h
    ) internal view returns (bool) {
        for (uint256 i = 0; i < k; ) {
            uint256 bitPos    = _getLocation(h, i, filterSizeBits);
            // Because filterSizeBits ≤ filter.length * 8, we know:
            //    bitPos < filterSizeBits ≤ filter.length*8
            // → (bitPos >> 3) < filter.length
            // Therefore the `byteIndex` check is redundant and can be skipped.
            uint256 byteIndex = bitPos >> 3;
            uint8  bitOffset  = uint8(bitPos & 7);

            // Read one byte from storage:
            uint8 b = uint8(filter[byteIndex]);
            if (((b >> bitOffset) & 1) == 0) {
                return false;
            }

            unchecked { ++i; }
        }
        return true;
    }

    /// @notice Extract four 64‐bit values from keccak256(token)
    function extractHashes(bytes calldata token) internal pure returns (uint64[4] memory h) {
        bytes32 digest = keccak256(token);
        h[0] = uint64(uint256(digest >> 192));
        h[1] = uint64((uint256(digest) >> 128) & 0xFFFFFFFFFFFFFFFF);
        h[2] = uint64((uint256(digest) >>  64) & 0xFFFFFFFFFFFFFFFF);
        h[3] = uint64(uint256(digest)          & 0xFFFFFFFFFFFFFFFF);
    }

    /// @notice “Go‐style” mixing of the four 64‐bit hashes into k positions
    function _getLocation(
        uint64[4] memory h,
        uint256 i,
        uint256 mod
    ) internal pure returns (uint256) {
        // We know i fits in uint64 because k ≤ type(uint32).max ≤ uint64
        uint64 ii   = uint64(i);
        uint64 base = h[ii & 1]; // same as h[ii % 2]
        // idx = 2 + ⌊((ii + (ii mod 2)) mod 4) / 2⌋
        //   Because (ii mod 2) is either 0 or 1, (ii + (ii mod 2)) mod 4 yields 0..3,
        //   dividing by 2 gives 0 or 1.  So idx ∈ {2,3}.
        uint64 idx  = 2 + uint64(((ii + (ii & 1)) & 3) >> 1);
        uint64 mult = h[idx];

        uint64 sum64;
        unchecked {
            sum64 = base + ii * mult;
        }
        return uint256(sum64) % mod;
    }
}