# neb
A web gateway to handle REST API requests and provide access to individual back end services, handle domain setup use cases, authentication, security, metrics, request/response transformations and logging

## Objectives

- Create a web api gateway for handling RESTful APIs and connect with services using gRPC and Protocol Buffer
- Inject services and repositories to keep them decoupled and define layers based on domain driven design concepts
- Setup a database connection pool with retries support and guard against sql injection
- Handle graceful server shutdown, provide middleware for logging and instrumentation of service functions
- Serve static pages and provide circuit breaker for services
- Handle access control, secure session cookies, secure headers, XSS, CSRF and input sanitization

## Structure

The code is based on the concepts defined in the Domain driven design approach. It has the following packages:
- server - handles routing, decoding requests, encoding responses, access control
- postgres - contains a postgresql implementation of the repository interface
- static - contains the static assets html, js, css
- authoring - contains the application service for the use case of documenting features and wireframes
- rde - pure domain package for the requirement development environment.

## Inversion of control

The application service and respository are defined as interfaces that allow the outer layers to inject concrete implementations, leaving the control of which implementation to use, to the caller.
