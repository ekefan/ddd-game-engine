# Project Folder Specification For a sample DDD application with  Hexagonal Architecture

project-root/
├── application/
│   ├── ports/
│   │   ├── inbound/
│   │   │   └── // Contains interfaces that define the entry points to the application. These ports are used by external agents interacting with the application, such as user interfaces or REST API requests.
│   │   └── outbound/
│   │       └── // Defines interfaces for external services that the application needs to consume, such as databases or external REST services. These ports help decouple the business logic from the implementation details of accessing external resources.
│   └── services/
│       └── // Implements the application logic coordinating activities between the ports and the domain. Application services play a crucial role in orchestrating domain operations, executing business logic, and acting as a bridge between the domain and infrastructure adapters.
├── domain/
│   ├── exceptions/
│   │   └── // Defines specific domain exceptions that can be thrown by the business logic.
│   ├── entities/
│   │   └── // Contains the domain entities that encapsulate critical business logic and data.
│   └── other domain folders/
│       └── // May include value objects, aggregates, domain events, etc., which are fundamental to the business logic and domain rules.
└── infrastructure/
    ├── adapters/
    │   ├── inbound/
    │   │   ├── rest/
    │   │   │   └── // Implements adapters for web interfaces, handling incoming HTTP requests and transforming them into calls to the appropriate inbound ports.
    │   │   ├── tasks/
    │   │   │   └── // For scheduled tasks that perform periodic operations within the application.
    │   │   └── events/
    │   │       └── // Manages the capture and processing of system events or integration events with other systems.
    │   └── outbound/
    │       ├── persistence/
    │       │   └── // Implements data persistence, for example, using JPA to interact with databases, encapsulating all data access logic.
    │       └── rest/
    │           └── // Contains the adapters needed to make calls to external APIs, encapsulating the logic of how to interact with other web services.
    └── configuration/
        └── // Specific configurations of the framework and infrastructure, such as security settings, Spring beans configuration, etc.
