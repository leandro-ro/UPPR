// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CascadingBloomFilter {
    /* ─── State ───────────────────────────────────────────────────────────── */

    address private _owner;

    struct Layer {
        bytes   filter;         // Packed bit‐vector
        uint256 filterSizeBits; // Number of valid bits in `filter`
        uint256 k;              // Number of hash functions
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
        require(newFilters.length == ks.length,  "Need k for each layer");
        require(ks.length == bitLens.length,     "Need bitLen for each layer");
        require(newFilters.length >  0,           "At least one layer");

        // Remove any existing layers
        delete layers;

        for (uint256 i = 0; i < newFilters.length; i++) {
            bytes    calldata f    = newFilters[i];
            uint256  k_            = ks[i];
            uint256  bits          = bitLens[i];

            require(k_ > 0,                          "k must be > 0");
            require(bits > 0,                        "bitLen must be > 0");
            require(bits <= f.length * 8,           "bitLen exceeds f.length * 8");

            // Create and populate a new Layer
            layers.push();
            Layer storage L = layers[layers.length - 1];
            L.filter         = f;
            L.filterSizeBits = bits;
            L.k              = k_;
        }
    }

    /// @notice Return the total number of layers.
    function layerCount() external view returns (uint256) {
        return layers.length;
    }

    /// @notice Test `token` against every layer. Returns (accepted, layerIndexReached).
    function testToken(bytes calldata token) external view returns (bool, uint256) {
        uint256 n = layers.length;
        require(n > 0, "No layers");

        uint64[4] memory h = extractHashes(token);

        for (uint256 li = 0; li < n; li++) {
            Layer storage L = layers[li];
            bool match_     = testInLayer(L.filter, L.filterSizeBits, L.k, h);

            // If this is the last layer, we expect a match on even‐indexed layers
            if (li == n - 1) {
                bool expected = (li % 2 == 0);
                return (match_ == expected, li);
            }

            // If there is a mismatch before the last layer, see if we accept early
            if (!match_) {
                bool acceptEarly = (li % 2 == 1);
                return (acceptEarly, li);
            }
        }

        revert("unreachable");
    }

    /// @notice Return metadata and full filter bytes for layer i.
    /// @param i Index of the layer (0‐based)
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
        return (L.filterSizeBits, L.k, L.filter);
    }

    /* ─── Internal Helpers ─────────────────────────────────────────────────── */

    /// @notice Test a single layer’s `filter` as a Bloom filter.
    /// @param filter         the full packed bit‐vector (in bytes)
    /// @param filterSizeBits number of bits of `filter` to consider
    /// @param k              number of hash probes
    /// @param h              4×64‐bit “precomputed” hashes
    function testInLayer(
        bytes storage filter,
        uint256 filterSizeBits,
        uint256 k,
        uint64[4] memory h
    ) internal view returns (bool) {
        for (uint256 i = 0; i < k; i++) {
            uint256 bitPos    = getLocation(h, i, filterSizeBits);
            uint256 byteIndex = bitPos >> 3;
            require(byteIndex < filter.length, "byteIndex out of range");

            uint8 bitOffset = uint8(bitPos & 7);
            uint8 b         = uint8(filter[byteIndex]);

            if (((b >> bitOffset) & 1) == 0) {
                return false;
            }
        }
        return true;
    }

    /// @notice Extract four 64‐bit values from keccak256(token)
    function extractHashes(bytes calldata token) internal pure returns (uint64[4] memory h) {
        bytes32 digest = keccak256(token);
        h[0] = uint64(uint256(digest >> 192));
        h[1] = uint64((uint256(digest) >> 128) & 0xFFFFFFFFFFFFFFFF);
        h[2] = uint64((uint256(digest) >>  64) & 0xFFFFFFFFFFFFFFFF);
        h[3] = uint64(uint256(digest)         & 0xFFFFFFFFFFFFFFFF);
    }

    /// @notice “Go‐style” mixing of the four 64‐bit hashes into k positions
    /// @param h   Array of 4 precomputed 64‐bit hashes
    /// @param i   Index of the hash probe (0 ≤ i < k)
    /// @param mod The Bloom filter’s bit size (filterSizeBits)
    function getLocation(
        uint64[4] memory h,
        uint256 i,
        uint256 mod
    ) internal pure returns (uint256) {
        uint64 ii   = uint64(i);
        uint64 base = h[ii % 2];
        // idx = 2 + ⌊((ii + (ii mod 2)) mod 4) / 2⌋
        uint64 idx  = 2 + (((ii + (ii % 2)) % 4) / 2);
        uint64 mult = h[idx];
        uint64 sum64;
        unchecked {
            sum64 = base + ii * mult;
        }
        return uint256(sum64) % mod;
    }
}