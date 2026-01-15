# Spec: Graceful Shutdown

## ADDED Requirements

### Requirement: Signal Handling
The system SHALL intercept termination signals for graceful shutdown.

#### Scenario: SIGTERM Signal
- **WHEN** the process receives SIGTERM
- **THEN** the system SHALL initiate graceful shutdown
- **AND** the system SHALL log "Received SIGTERM, initiating graceful shutdown"
- **AND** the system SHALL not accept new gRPC connections
- **AND** the system SHALL begin draining active streams

#### Scenario: SIGINT Signal
- **WHEN** the process receives SIGINT (Ctrl+C)
- **THEN** the system SHALL initiate graceful shutdown
- **AND** the system SHALL log "Received SIGINT, initiating graceful shutdown"
- **AND** the system SHALL follow the same shutdown sequence as SIGTERM

#### Scenario: Multiple Signals
- **WHEN** multiple termination signals are received
- **THEN** the system SHALL ignore subsequent signals
- **AND** the system SHALL complete the first shutdown sequence
- **AND** if shutdown exceeds timeout, the system SHALL forcefully terminate

### Requirement: Binance Connection Closure
The system SHALL cleanly close the Binance WebSocket connection during shutdown.

#### Scenario: Binance Connection Close
- **WHEN** graceful shutdown is initiated
- **THEN** the system SHALL send a close frame to Binance WebSocket
- **AND** the system SHALL wait up to 2 seconds for close acknowledgment
- **AND** the system SHALL log "Closed Binance WebSocket connection"

#### Scenario: Binance Close Timeout
- **WHEN** the Binance WebSocket does not acknowledge the close
- **THEN** the system SHALL forcibly close the connection after 2 seconds
- **AND** the system SHALL log "Forcibly closed Binance WebSocket connection"

### Requirement: Stream Draining
The system SHALL drain active gRPC streams before termination.

#### Scenario: Drain Window
- **WHEN** graceful shutdown is initiated
- **THEN** the system SHALL allow 5 seconds for stream draining
- **AND** the system SHALL send a `ShutdownNotification` to all active streams
- **AND** the system SHALL continue processing updates during the drain window

#### Scenario: Client Notification
- **WHEN** shutdown is initiated
- **THEN** each connected client SHALL receive a final message
- **AND** the message SHALL indicate the server is shutting down
- **AND** the message SHALL allow clients to gracefully close their side

#### Scenario: Stream Completion Within Drain
- **WHEN** all active streams complete before the 5-second timeout
- **THEN** the system SHALL proceed to server shutdown immediately
- **AND** the system SHALL log "All streams drained successfully"

### Requirement: gRPC Server Shutdown
The system SHALL shutdown the gRPC server gracefully.

#### Scenario: GracefulStop Call
- **WHEN** stream draining is complete (or timeout reached)
- **THEN** the system SHALL call `grpc.Server.GracefulStop()`
- **AND** the system SHALL wait up to 5 seconds for in-flight requests to complete
- **AND** the system SHALL log "gRPC server graceful stop initiated"

#### Scenario: GracefulStop Timeout
- **WHEN** `GracefulStop()` does not complete within 5 seconds
- **THEN** the system SHALL call `grpc.Server.Stop()`
- **AND** the system SHALL forcibly terminate all connections
- **AND** the system SHALL log "gRPC server forced stop due to timeout"

### Requirement: Resource Cleanup
The system SHALL release all resources during shutdown.

#### Scenario: Memory Cleanup
- **WHEN** shutdown is complete
- **THEN** the system SHALL clear all price history buffers
- **AND** the system SHALL close all internal channels
- **AND** the system SHALL release all goroutines

#### Scenario: Logging Final State
- **WHEN** shutdown is complete
- **THEN** the system SHALL log "Graceful shutdown completed successfully"
- **OR** the system SHALL log "Shutdown completed with errors: <details>" if issues occurred

### Requirement: Shutdown Sequence Summary
The complete shutdown sequence SHALL execute in this order.

#### Scenario: Full Shutdown Timeline
- **GIVEN** the server is running with active client connections
- **WHEN** SIGTERM is received
- **THEN** the shutdown sequence SHALL complete within 10 seconds maximum
- **AND** the server SHALL exit with code 0 on success
