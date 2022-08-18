# starbridge-client

This project was bootstrapped with [Create React App](https://create-react-app.dev/) using the [Typescript template](https://github.com/facebook/create-react-app/tree/master/packages/cra-template-typescript).

## Project Architecture

You can check the project architecture [here](./src/docs/ARCHITECTURE.md)

## Requirements

- [NodeJS 14+](https://nodejs.org/en/)
- [NPM 6.14+](https://www.npmjs.com/)

> We suggest use of [NVM](https://github.com/nvm-sh/nvm/blob/master/README.md) to manage your node versions.

## Getting Started

### Env vars config

The environment variables are in `src/config`. You can use the `.env.example` as a base to create your `.env.local`
config file

### Install dependencies

```shell
npm install
```

### Running in development environment

```shell
npm run start:dev
```

The project will be running at [http://localhost:3000/](http://localhost:3000/)

### Running tests

- We use the [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/) to develop our tests.
- You can use the [MSW](https://mswjs.io/) to mock your request to do integration tests in your pages.
  - An example is available at `src/tests/request_mocks`

```shell
npm run test
```

### Creating a production build

The following command will generate an optimized production build. The statics files will be generated at `build/` folder.

```shell
npm run build
```

You can read more about how to serve the statics [here](https://create-react-app.dev/docs/deployment/)

## Scripts

In the project directory, you can run all of [react-scripts](https://create-react-app.dev/docs/available-scripts) commands.

## Learn More

To learn React, check out the [React documentation](https://reactjs.org/).
