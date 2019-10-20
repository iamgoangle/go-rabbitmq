# GO AMQP Client Wrapper

## About

Implementation is probably straight-forward of the project that aim to follow chain-of-responsiblity and [functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis) pattern to solve the complexity and make abstraction layer to a user that use the this library.

### Features

- [x] Middleware pattern for initialization producer.

- [x] Functional options for producer publishing the message.

- [ ] Connection re-establish once recieves closed notifty signal from broker as soon as possible.

- [ ] Consumer retry queue middleware.

- [ ] Consumer retry queue middleware with backoff config.

- [x] Closed Connection and Channel gracefully by signal notify

### Usages
