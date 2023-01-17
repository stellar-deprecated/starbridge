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
    pub fn init(env: Env, admin: Identifier) {
        panic!("unimplemented");
    }

    pub fn deposit(env: Env, token: BytesN<32>, is_wrapped_asset: bool, eth_destination: AccountId, amount: i128) {
        let client = token::Client::new(&env, &token);
        let from = &env.invoker().into();

        if is_wrapped_asset {
            client.burn_from(&Signature::Invoker, &0, &from, &amount);
        } else {
            client.xfer_from(&Signature::Invoker, &0, from, &Identifier::Contract(env.current_contract()), &amount);
        }

        let topics = (symbol!("deposit"), &token, from, eth_destination);
        env.events().publish(topics, amount);
    }

    pub fn withdraw(env: Env, token: BytesN<32>, is_wrapped_asset: bool, recipient: Identifier, id: BytesN<32>, expiration: u64) {
        panic!("unimplemented");
    }

    pub fn set_paused(env: Env, state : Pause) {
        panic!("unimplemented");
    }

    pub fn fulfilled(env: Env, id: BytesN<32>) -> bool {
        panic!("unimplemented");
    }

    pub fn admin(env: Env) -> Identifier {
        panic!("unimplemented");
    }
}
