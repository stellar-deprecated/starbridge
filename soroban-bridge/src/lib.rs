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
        if env.storage().has(&key) {
            panic!("admin already initialized!");
        }
        env.storage().set(&key, admin);
    }

    pub fn deposit(env: Env, token: BytesN<32>, is_wrapped_asset: bool, eth_destination: AccountId, amount: i128) {
        if !env.storage().has(DataKey::Admin) {
            panic!("contract not initialized!");
        }
        if get_paused(&env).map(|p| p == Pause::All || p == Pause::Deposit).unwrap_or(false) {
            panic!("deposits are paused!")
        }
        if amount <= 0 {
            panic!("only positive amounts allowed!")
        }

        let client = token::Client::new(&env, &token);
        let from = &env.invoker().into();

        if is_wrapped_asset {
            client.burn_from(&Signature::Invoker, &0, from, &amount);
        } else {
            client.xfer_from(&Signature::Invoker, &0, from, &Identifier::Contract(env.current_contract()), &amount);
        }
        
        let topics = (symbol!("deposit"), &token, from, eth_destination);
        env.events().publish(topics, amount);
    }

    pub fn withdraw(env: Env, token: BytesN<32>, is_wrapped_asset: bool, recipient: Identifier, id: BytesN<32>, amount: i128) {
        if get_paused(&env).map(|p| p == Pause::All || p == Pause::Withdrawal).unwrap_or(false) {
            panic!("withdrawals are paused!")
        }
        if amount <= 0 {
            panic!("only positive amounts allowed!")
        }
        check_admin(&env, &env.invoker().into());

        let topics = (symbol!("withdraw"), &id, &token, &recipient);
        env.events().publish(topics, amount);

        let key = DataKey::Fullfilled(id);
        if env.storage().has(&key) {
            panic!("withdrawal already fulfilled");
        } else {
            env.storage().set(&key, ());
        }

        let client = token::Client::new(&env, &token);
        if is_wrapped_asset {
            client.mint(&Signature::Invoker, &0, &recipient, &amount);
        } else {
            client.xfer(&Signature::Invoker, &0, &recipient, &amount)
        }
    }

    pub fn set_paused(env: Env, state : Option<Pause>) {
        check_admin(&env, &env.invoker().into());
        let key = DataKey::Pause;
        match state {
            None => env.storage().remove(key),
            Some(p) => env.storage().set(key, p),
        }
    }

    pub fn status(env: Env, id: BytesN<32>) -> (bool, u32, u64) {
        let fulfilled = env.storage().has(DataKey::Fullfilled(id));
        let seq = env.ledger().sequence();
        let timestamp = env.ledger().timestamp();
        (fulfilled, seq, timestamp)
    }

    pub fn admin(env: Env) -> Identifier {
        get_admin(&env)
    }
}

pub fn get_paused(e: &Env) -> Option<Pause> {
    e.storage().get(DataKey::Pause).map(|p| p.unwrap())
}

fn get_admin(e: &Env) -> Identifier {
    let key = DataKey::Admin;
    e.storage().get_unchecked(key).unwrap()
}

pub fn check_admin(e: &Env, auth_id: &Identifier) {
    if *auth_id != get_admin(e) {
        panic!("not authorized by admin")
    }
}