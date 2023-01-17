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
#[derive(PartialEq)]
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
        let key = DataKey::Admin;
        if env.storage().has(&key){
            panic!("admin already initialized!");
        }
        env.storage().set(&key, admin);
        env.storage().set(DataKey::Pause, Pause::None);
    }

    pub fn deposit(env: Env, token: BytesN<32>, is_wrapped_asset: bool, eth_destination: AccountId, amount: i128) {
        let paused: Pause = env.storage().get_unchecked(DataKey::Pause).unwrap();
        if paused == Pause::All || paused == Pause::Deposit{
            panic!("deposits are paused!")
        }
        if amount < 0 {
            panic!("negative amount is not allowed!")
        }

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

    pub fn withdraw(env: Env, token: BytesN<32>, is_wrapped_asset: bool, recipient: Identifier, id: BytesN<32>) {
        panic!("unimplemented");
    }

    pub fn set_paused(env: Env, state : Pause) {
        check_admin(&env, &env.invoker().into());
        env.storage().set(DataKey::Pause, state);
    }

    pub fn fulfilled(env: Env, id: BytesN<32>) -> bool {
        panic!("unimplemented");
    }

    pub fn admin(env: Env) -> Identifier {
        read_administrator(&env)
    }
}

fn read_administrator(e: &Env) -> Identifier {
    let key = DataKey::Admin;
    e.storage().get_unchecked(key).unwrap()
}

pub fn check_admin(e: &Env, auth_id: &Identifier) {
    if *auth_id != read_administrator(e) {
        panic!("not authorized by admin")
    }
}