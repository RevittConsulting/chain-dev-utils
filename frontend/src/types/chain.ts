export type Chain = {
  networkName: string;
  L1: L1;
  L2: L2;
  lastUpdated: Date;
  serviceName: string;
};

type L1 = {
  chainId: string;
  rpcUrl: string;

  etherscan: string;
  blockscout: string;

  rollupManagerAddress: string;
  rollupAddress: string;

  latestL1BlockNumber: number;
  highestSequencedBatch: number;
  highestVerifiedBatch: number;
};

type L2 = {
  chainId: string;
  rpcUrl: string;
  polygonscan: string;
  datastreamerUrl: string;

  latestBatchNumber: number;
  latestBlockNumber: number;
  datastreamerStatus: string;
};

export type ChainData = {
  mostRecentL1Block: number;
  highestSequencedBatch: number;
  highestVerifiedBatch: number;
  mostRecentL2Batch: number;
  mostRecentL2Block: number;
  dataStreamerStatus: string;
}