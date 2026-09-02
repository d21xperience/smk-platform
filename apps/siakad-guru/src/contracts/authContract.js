export const AuthContract = {
  login: {
    request: {
      username: 'string',
      password: 'string',
    },
    response: {
      token: 'string',
      expiresAt: 'string', // ISO date format
      user: {
        id: 'string',
        name: 'string',
        email: 'string',
        role: 'string',
        permissions: ['string'],
      },
    },
  },
  logout: {
    request: {},
    response: {
      success: 'boolean',
    },
  },
}
