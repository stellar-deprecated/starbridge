#![no_std]

use soroban_auth::{Identifier, Signature};
use soroban_sdk::{contractimpl, symbol, contracttype, AccountId, BytesN, Env};

mod token {
    soroban_sdk::contractimport!(file = "./soroban_token_spec.wasm");
}

#[derive(Clone)]
#[contracttype]
pub enum DataKey {
    Admin,
    Fullfilled(BytesN<32>),
    Pause,
}

#[derive(Clone)]
#[contracttype]
pub enum Pause {
    None,
    Deposit,
    Withdrawal,
    All,
}

pub struct Bridge;

#[contractimpl]
#[allow(unused_variables)]
impl Bridge {
    pub fn init(env: Env, admin: Signature) {
        panic!("unimplemented");
    }

    pub fn deposit(env: Env, from: Signature, token: BytesN<32>, eth_destination: AccountId, amount: i128) {
        let topics = (symbol!("deposit"), from.identifier(&env), eth_destination);
        env.events().publish(topics, amount);
        panic!("unimplemented");
    }

    pub fn withdraw(env: Env, admin: Signature, token: BytesN<32>, recipient: Identifier, id: BytesN<32>, expiration: u64) {
        panic!("unimplemented");
    }

    pub fn set_paused(env: Env, admin: Signature, state : Pause) {
        panic!("unimplemented");
    }

    pub fn fulfilled(env: Env, id: BytesN<32>) -> bool {
        panic!("unimplemented");
    }

    pub fn admin(env: Env) -> Identifier {
        panic!("unimplemented");
    }
}
