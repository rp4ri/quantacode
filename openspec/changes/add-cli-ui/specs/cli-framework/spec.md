## ADDED Requirements

### Requirement: CLI Framework Initialization
The system SHALL initialize a CLI application using Cobra framework with a root command and 'chat' subcommand.

#### Scenario: Initialize root command
- **WHEN** CLI application starts
- **THEN** root command is available with help text

#### Scenario: Initialize chat subcommand
- **WHEN** user runs `quantacode chat`
- **THEN** TUI interface is launched

### Requirement: Server Connection Flag
The system SHALL provide a --server flag to specify gRPC server address with default value localhost:50051.

#### Scenario: Use default server address
- **WHEN** CLI starts without --server flag
- **THEN** it connects to localhost:50051

#### Scenario: Use custom server address
- **WHEN** CLI starts with `--server=example.com:50051`
- **THEN** it connects to specified address

### Requirement: Connection Status Display
The system SHALL display connection status to the gRPC server on startup.

#### Scenario: Show successful connection
- **WHEN** gRPC connection succeeds
- **THEN** status message shows "Connected to server at <address>"

#### Scenario: Show connection error
- **WHEN** gRPC connection fails
- **THEN** error message shows reason and offers options

### Requirement: Graceful Error Handling
The system SHALL handle CLI startup errors gracefully without crashing.

#### Scenario: Handle missing dependencies
- **WHEN** required flags are invalid or missing
- **THEN** CLI shows error message and usage instructions
