#!make


help:
	@echo "Usage: make [build|start|start_dev|test|storybook|build_storybook|lint|format_code"
	@echo ""
	@echo "Usage:"
	@echo "  make build  	        Builds the frontend for production and copies it to the dist folder"
	@echo "  make start  	        Starts the frontend in production mode"
	@echo "  make start_dev         Starts the frontend in development mode"
	@echo "  make test  	        Runs the frontend tests"
	@echo "  make coverage-test     Run tests on the project and generates a coverage report."
	@echo "  make storybook         Starts the storybook"
	@echo "  make build_storybook   Builds the storybook"
	@echo "  make lint  	        Runs the frontend linter"
	@echo "  make format_code       Formats the frontend code with prettier and fix the code style"
	@echo ""

setup: package.json
	npm install

build: setup
	npm run build

start: setup
	npm run start

start_dev: setup
	npm run start:dev

test: setup
	npm run test

coverage-test: setup
	npm run test -- --coverage

storybook: setup
	npm run storybook

build-storybook: setup
	npm run build-storybook

lint: setup
	npm run lint

format-code: setup
	npm run prettier
