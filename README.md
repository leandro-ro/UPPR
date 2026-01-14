# UPPR: Universal Privacy-Preserving Revocation
This repository is the proof-of-concept implementation of **UPPR**, a universal privacy-preserving revocation framework for Verifiable Credentials (VCs) presented at the IEEE International Conference on Blockchain 2025 (https://ieeexplore.ieee.org/document/11264637). UPPR supports both linkable (one-show/oVC) and unlinkable (multi-show/AC) credentials by combining Verifiable Random Functions (VRFs) with a scalable Bloom filter cascade. It enables efficient, metadata-free revocation without requiring holders to interact with issuers during credential presentation.

## Project Structure

### `bloom`
Implements the Bloom filter cascade used for encoding revocation artifacts.
- `sol/`: Solidity implementation for on-chain verification.
- `cascade.go`: Go implementation for off-chain artifact construction.
- `filter.go`: Bloom filter logic adapted from [bits-and-blooms/bloom](https://github.com/bits-and-blooms/bloom/blob/master/bloom.go).

### `external`
Contains external dependencies and adapted libraries.
- `go-ecvrf/`: Fork of [vechain/go-ecvrf](https://github.com/vechain/go-ecvrf) with improved EC operations using [go-ethereum](https://github.com/ethereum/go-ethereum).
- `vrf/`: Solidity contracts from [witnet/vrf-solidity](https://github.com/witnet/vrf-solidity) for on-chain VRF functionality.

### `holder`
Implements holder-side logic for generating non-revocation proofs using credentials and revocation artifacts.

### `issuer`
Implements issuer-side logic for credential issuance and revocation artifact generation.

### `verifier`
Contains two Solidity smart contracts:
- A verifier for one-show credentials (oVC).
- A verifier for multi-show credentials (AC) using Zero-Knowledge Proofs.

### `zkp`
Implements the Zero-Knowledge circuit for multi-show credential revocation using [gnark](https://github.com/Consensys/gnark) and provides the corresponding Solidity verifier for on-chain validation.

## Usage
All packages in this repository include comprehensive tests. You can run the full test suite from the root directory with:

    go test ./...

To execute all available benchmarks, use:

    go test -bench . -run=^$ -benchmem ./...

This will run all benchmarks across the codebase and display performance and memory statistics.

---

## Benchmarks

We provide comprehensive benchmarks to assess runtime and on-chain performance of key operations in UPPR. Below is a summary of selected results.

### Revocation Artifact Generation

Artifact generation scales linearly with the number of issued credentials. Multi-show credentials are significantly cheaper to process than one-show credentials due to the use of hash functions instead of VRFs.

| Credential Type | Domain | Revocation Rate | N | Time per Operation (ns/op) |
| --- | --- | --- | --- | --- |
| OneShow | 50000 | 0.10 | 10 | 209533808 |
| OneShow | 100000 | 0.10 | 10 | 417011904 |
| OneShow | 200000 | 0.10 | 10 | 836083775 |
| OneShow | 300000 | 0.10 | 10 | 1271718658 |
| OneShow | 400000 | 0.10 | 10 | 1721598354 |
| OneShow | 500000 | 0.10 | 10 | 2095792862 |
| OneShow | 600000 | 0.10 | 10 | 2553666121 |
| OneShow | 700000 | 0.10 | 10 | 2965309096 |
| OneShow | 800000 | 0.10 | 10 | 3417696242 |
| OneShow | 900000 | 0.10 | 10 | 3842405254 |
| OneShow | 1000000 | 0.10 | 10 | 4231854058 |
| MultiShow | 50000 | 0.10 | 10 | 104810962 |
| MultiShow | 100000 | 0.10 | 10 | 182449854 |
| MultiShow | 200000 | 0.10 | 10 | 366856717 |
| MultiShow | 300000 | 0.10 | 10 | 547776792 |
| MultiShow | 400000 | 0.10 | 10 | 729738012 |
| MultiShow | 500000 | 0.10 | 10 | 920114858 |
| MultiShow | 600000 | 0.10 | 10 | 1098911667 |
| MultiShow | 700000 | 0.10 | 10 | 1279644150 |
| MultiShow | 800000 | 0.10 | 10 | 1517099583 |
| MultiShow | 900000 | 0.10 | 10 | 1712877912 |
| MultiShow | 1000000 | 0.10 | 10 | 1944766500 |
### Holder Proof Generation

| Type      | Main Primitive | Time [ms] | Memory [MB] | On-Chain Verification [gas] |
|-----------|---------------|-----------|-------------|-----------------------------|
| OneShow   | VRF.Eval      | 0.072     | 0.004       | 422209                      |
| MultiShow | ZKP.Prove     | 33.74     | 5.34        | 281741                      |

### On-Chain Update Cost

The table below shows the gas and ETH cost for updating the Bloom filter cascade on-chain. The first update writes the filter to storage, while the second only modifies existing entries, making it significantly cheaper.

| Domain   | Capacity | 1st Avg Gas | 1st ETH        | 2nd Avg Gas | 2nd ETH        |
|----------|----------|-------------|----------------|-------------|----------------|
| 50000    | 5%       | 3,035,721   | 0.003035721 ETH| 853,150     | 0.000853150 ETH|
| 50000    | 10%      | 4,742,556   | 0.004742556 ETH| 1,376,237   | 0.001376237 ETH|
| 200000   | 5%       | 10,154,585  | 0.010154585 ETH| 2,837,370   | 0.002837370 ETH|
| 200000   | 10%      | 16,677,645  | 0.016677645 ETH| 4,636,050   | 0.004636050 ETH|
| 400000   | 5%       | 19,446,435  | 0.019446435 ETH| 5,367,919   | 0.005367919 ETH|
| 400000   | 10%      | 32,448,724  | 0.032448724 ETH| 8,881,321   | 0.008881321 ETH|
| 600000   | 5%       | 28,708,334  | 0.028708334 ETH| 7,878,829   | 0.007878829 ETH|
| 600000   | 10%      | 48,067,620  | 0.048067620 ETH| 13,172,435  | 0.013172435 ETH|
| 800000   | 5%       | 38,012,533  | 0.038012533 ETH| 10,376,464  | 0.010376464 ETH|
| 800000   | 10%      | 63,809,878  | 0.063809878 ETH| 17,375,514  | 0.017375514 ETH|
| 1000000  | 5%       | 47,056,090  | 0.047056090 ETH| 12,876,066  | 0.012876066 ETH|
| 1000000  | 10%      | 79,404,093  | 0.079404093 ETH| 21,666,129  | 0.021666129 ETH|

### End-to-End One-Show Verification

Benchmark gas consumption for verifying a one-show credential presentation using `CheckCredential` (N = 500 credentials):

| Domain   | Capacity | Avg Gas Used | ETH (1 Gwei) | Local Check [ms] |
|----------|----------|---------------|--------------|------------------|
| 1000     | 5%       | 422209        | 0.000422209  | 0.745            |
| 1000     | 10%      | 417678        | 0.000417678  | 0.749            |
| 10000    | 5%       | 422580        | 0.000422580  | 0.774            |
| 10000    | 10%      | 428846        | 0.000428846  | 0.814            |
| 100000   | 5%       | 419191        | 0.000419191  | 0.794            |
| 100000   | 10%      | 430095        | 0.000430095  | 0.766            |
| 1000000  | 5%       | 423070        | 0.000423070  | 0.756            |
| 1000000  | 10%      | 426334        | 0.000426334  | 0.801            |

### End-to-End Multi-Show Verification

Benchmark gas consumption for verifying a multi-show credential presentation using `CheckCredential` (N = 500 credentials):

| Domain   | Capacity | Avg Gas Used | ETH (1 Gwei) | Local Check [ms] |
|----------|----------|--------------|--------------|------------------|
| 1000     | 5%       | 279844       | 0.000279844  | 2.097            |
| 1000     | 10%      | 280636       | 0.000280636  | 2.153            |
| 10000    | 5%       | 281258       | 0.000281258  | 2.147            |
| 10000    | 10%      | 282025       | 0.000282025  | 2.179            |
| 100000   | 5%       | 282119       | 0.000282119  | 2.172            |
| 100000   | 10%      | 283309       | 0.000283309  | 2.171            |
| 1000000  | 5%       | 281741       | 0.000281741  | 2.162            |
| 1000000  | 10%      | 281661       | 0.000281661  | 2.173            |

## Citation
If you use this repository or build upon UPPR, please cite the following paper:

```bibtex
@INPROCEEDINGS{11264637,
  author    = {Rometsch, Leandro and Lehwalder, Philipp-Florens and Hoang, Anh-Tu and Kaaser, Dominik and Schulte, Stefan},
  title     = {{UPPR: Universal Privacy-Preserving Revocation}},
  booktitle = {{2025 IEEE International Conference on Blockchain (Blockchain)}},
  year      = {2025},
  pages     = {161--170},
  doi       = {10.1109/Blockchain67634.2025.00030}
}
```

## Acknowledgment
The financial support by the Austrian Federal Ministry of Economy, Energy and Tourism, the National Foundation for Research, Technology and Development and the Christian Doppler Research Association is gratefully acknowledged. Further, this result is part of a project that received funding from the European Research Council (ERC) under the European Union’s Horizon 2020 and Horizon Europe research and innovation programs (grant CRYPTOLAYER-101044770). We also thank Mirko Mollik of the German Federal Agency for Breakthrough Innovation (SPRIND) for providing valuable industry insights.

## Copyright
Use of the source code is governed by the Apache 2.0 license that can be found in the [LICENSE file](LICENSE.txt).
