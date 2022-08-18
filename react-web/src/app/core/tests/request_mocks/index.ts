import { rest } from 'msw'

const BASE_HOST_URL = '*'

const handlers = [
  // The code above is an example about how you can use MSW to mock your requests
  rest.get(`${BASE_HOST_URL}/ping`, (req, res, ctx) => {
    return res(
      // Respond with a 200 status code
      ctx.status(200, 'pong')
    )
  }),
]

export { handlers }
