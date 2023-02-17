import { Injectable } from '@nestjs/common';
import {
  deserializeReceiveReturnValue,
  HttpProvider,
  JsonRpcClient,
  serializeUpdateContractParameters,
} from '@concordium/common-sdk';
import { Buffer } from 'buffer/';
import { RequestDto, TransactionDto } from './serialize-params.dto';
import {
  BRIDGE_CONTRACT_RAW_SCHEMA,
  TOKEN_CONTRACT_RAW_SCHEMA,
} from './consts';

@Injectable()
export class AppService {
  async getWithdrawHash(request: RequestDto): Promise<string> {
    const param = serializeUpdateContractParameters(
      'gbm_Bridge',
      'withdraw_hash',
      request.parameters,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
    );
    const gRPCProvider = new HttpProvider(
      'https://concordium-json-rpc.bridge.bankofmemories.org',
    );
    const rpcClient = new JsonRpcClient(gRPCProvider);
    const res = await rpcClient.invokeContract({
      method: 'gbm_Bridge.withdraw_hash',
      contract: {
        index: 2945n,
        subindex: 0n,
      },
      parameter: param,
    });
    if (!res || res.tag === 'failure' || !res.returnValue) {
      throw new Error(
        `RPC call 'invokeContract' on method gbm_Bridge.withdraw_hash of contract 2945 failed - ${res}`,
      );
    }
    const returnValues = deserializeReceiveReturnValue(
      Buffer.from(res.returnValue, 'hex') as any,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
      'gbm_Bridge',
      'withdraw_hash',
      0,
    );
    // eslint-disable-next-line no-console
    console.log(returnValues);

    return returnValues;
  }

  async withdraw(request: RequestDto): Promise<string> {
    const param = serializeUpdateContractParameters(
      'gbm_Bridge',
      'withdraw',
      request.parameters,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
    );
    const gRPCProvider = new HttpProvider(
      'https://concordium-json-rpc.bridge.bankofmemories.org',
    );
    const rpcClient = new JsonRpcClient(gRPCProvider);
    console.log(param.toString());
    const res = await rpcClient.invokeContract({
      method: 'gbm_Bridge.withdraw',
      contract: {
        index: 2945n,
        subindex: 0n,
      },
      parameter: param,
    });
    if (!res || res.tag === 'failure' || !res.returnValue) {
      throw new Error(
        `RPC call 'invokeContract' on method gbm_Bridge.withdraw of contract 2945 failed`,
      );
    }
    const returnValues = deserializeReceiveReturnValue(
      Buffer.from(res.returnValue, 'hex') as any,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
      'gbm_Bridge',
      'withdraw',
      0,
    );
    // eslint-disable-next-line no-console
    console.log(returnValues);

    return returnValues;
  }

  async deposit(request: RequestDto): Promise<string> {
    const param = serializeUpdateContractParameters(
      'gbm_Bridge',
      'deposit',
      request.parameters,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
    );
    const gRPCProvider = new HttpProvider(
      'https://concordium-json-rpc.bridge.bankofmemories.org',
    );
    const rpcClient = new JsonRpcClient(gRPCProvider);
    const res = await rpcClient.invokeContract({
      method: 'gbm_Bridge.deposit',
      contract: {
        index: 2945n,
        subindex: 0n,
      },
      parameter: param,
    });
    if (!res || res.tag === 'failure' || !res.returnValue) {
      throw new Error(
        `RPC call 'invokeContract' on method gbm_Bridge.deposit of contract 2945 failed`,
      );
    }
    const returnValues = deserializeReceiveReturnValue(
      Buffer.from(res.returnValue, 'hex') as any,
      Buffer.from(BRIDGE_CONTRACT_RAW_SCHEMA, 'base64') as any,
      'gbm_Bridge',
      'deposit',
      0,
    );
    // eslint-disable-next-line no-console
    console.log(returnValues);

    return returnValues;
  }

  async getDepositParams(request: TransactionDto): Promise<object> {
    const gRPCProvider = new HttpProvider(
      'https://concordium-json-rpc.bridge.bankofmemories.org',
    );
    const rpcClient = new JsonRpcClient(gRPCProvider);
    async function getTransactionStatus() {
      return await new Promise(function (resolve) {
        setTimeout(
          () => resolve(rpcClient.getTransactionStatus(request.hash)),
          1000,
        );
      });
    }
    let res;
    do {
      res = await getTransactionStatus();
      if (!res) {
        throw new Error(`RPC call 'getTransactionStatus' failed`);
      }
    } while ((await res).status !== 'finalized');
    const blockHash = Object.keys(res.outcomes)[0];
    const event = res.outcomes[blockHash].result['events'].find(
      (result) => result.receiveName === 'gbm_Bridge.deposit',
    );
    const message = event.message;
    const from = event.instigator.address;
    const serializedTransaction = Buffer.from(message, 'hex');
    const serializedDestination = serializedTransaction.slice(200, 256);
    const serializedAmount = serializedTransaction.slice(
      256,
      serializedTransaction.length,
    );
    return {
      amount: serializedAmount.readBigUInt64LE(0).toString(),
      destination: Buffer.from(serializedDestination).toString(),
      blockHash,
      from,
    };
  }

  async getBalanceOf(request: RequestDto): Promise<string> {
    const param = serializeUpdateContractParameters(
      'wGBM',
      'balanceOf',
      request.parameters,
      Buffer.from(TOKEN_CONTRACT_RAW_SCHEMA, 'base64') as any,
    );
    const gRPCProvider = new HttpProvider(
      'https://concordium-json-rpc.bridge.bankofmemories.org',
    );
    const rpcClient = new JsonRpcClient(gRPCProvider);
    const res = await rpcClient.invokeContract({
      method: 'wGBM.balanceOf',
      contract: {
        index: 2928n,
        subindex: 0n,
      },
      parameter: param,
    });
    if (!res || res.tag === 'failure' || !res.returnValue) {
      throw new Error(
        `RPC call 'invokeContract' on method gbm_Bridge.balanceOf of contract 2945 failed`,
      );
    }
    const returnValues = deserializeReceiveReturnValue(
      Buffer.from(res.returnValue, 'hex') as any,
      Buffer.from(TOKEN_CONTRACT_RAW_SCHEMA, 'base64') as any,
      'wGBM',
      'balanceOf',
      0,
    );
    // eslint-disable-next-line no-console
    console.log(returnValues);

    return returnValues;
  }
}
