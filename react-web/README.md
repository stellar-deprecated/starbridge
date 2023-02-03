
# Starbridge - Frontend

The objective of this project is to transfer assets between wallets from different networks.
Currently the available flows are from **Stellar to Ethereum** and vice versa, being able to send only **ETH** from one to the other.
With some initial setup you will be able to connect your wallets using **Freighter** for Stellar and **Metamask** for Ethereum and start the transfer flow.

⚠️ Remembering that the **refund** flow was not implemented in both transfer flows!

## Project Architecture

**Front:** React, Typescript using the Atomic Design

You can check the project architecture [here](./src/docs/ARCHITECTURE.md)

## Environment Variables

The environment variables are in `src/config`. You can use the `.env.example` as a base to create your `.env.local`
config file


## Run Locally

Install dependencies

```bash
  npm install
```

Start the server

```bash
  npm run start:dev
```

The project will be running at [http://localhost:3000/](http://localhost:3000/)

If you want you can run it on **Docker** too.

```bash
  docker build -t starbridge-front .   
  docker run -dp 3000:3000 starbridge-front
```

### Using the ngrok

To make transactions on your wallets, you must have a secure connection locally. With [ngrok](https://ngrok.com/download) it is possible to create an SSL certificate to be able to continue in the flow without problems.

You need to [register](https://dashboard.ngrok.com/signup), to generate a token and then continue configuring ngrok.

```bash
  ngrok config add-authtoken <token>
  ngrok http 3000
```

## FAQ

#### I disconnected my account through the page, but when I connect again it automatically connects.

It is necessary that you disconnect from the page through your Wallet as well, so you can connect another account again.

#### I'm trying to do the withdraw action on the Ethereum -> Stellar flow but the error "retry later once the transaction has more confirmations" appears.

Confirmations are required to verify and legitimize information that will be recorded in the blockchain and cannot be changed afterward. If some information is assumed fraudulent, it will not get any confirmation. Without a single transaction confirmation Ethereum, the transaction won’t be considered valid by the network. Each confirmation takes less than one minute. Just wait a bit and try to perform the withdraw action again.



## Demo

![](https://s4.gifyu.com/images/ezgif.com-gif-maker-1940da4e499b61722.gif)
