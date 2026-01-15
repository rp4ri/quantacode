## ADDED Requirements

### Requirement: Bubble Tea TUI Framework
The system SHALL implement a TUI using Bubble Tea framework with Model-View-Update architecture.

#### Scenario: Initialize Bubble Tea program
- **WHEN** chat subcommand is executed
- **THEN** Bubble Tea program starts with initial model

#### Scenario: Implement Update method
- **WHEN** messages are received
- **THEN** model updates state and returns commands

#### Scenario: Implement View method
- **WHEN** model state changes
- **THEN** view renders current state to terminal

### Requirement: Split Screen Layout
The system SHALL provide a split-screen layout with price ticker on top and chat messages below.

#### Scenario: Display price ticker
- **WHEN** price update is received
- **THEN** ticker shows symbol, current price, and percent change

#### Scenario: Display chat messages
- **WHEN** user or AI sends a message
- **THEN** message appears in chat area with role indicator

### Requirement: Price Change Styling
The system SHALL use Lipgloss to color-code price changes: green for increase, red for decrease.

#### Scenario: Green color for price increase
- **WHEN** current price > previous price
- **THEN** price change percentage displays in green

#### Scenario: Red color for price decrease
- **WHEN** current price < previous price
- **THEN** price change percentage displays in red

### Requirement: User Input Textarea
The system SHALL provide a textarea at bottom for user to type questions.

#### Scenario: Display input prompt
- **WHEN** TUI is running
- **THEN** input area shows "Type your question..." prompt

#### Scenario: Accept user input
- **WHEN** user types text
- **THEN** characters appear in textarea

### Requirement: Question Submission
The system SHALL send user questions to gRPC server when Enter key is pressed.

#### Scenario: Send question on Enter
- **WHEN** user types question and presses Enter
- **THEN** question is sent to gRPC StreamAIAnalysis method

#### Scenario: Clear input after send
- **WHEN** question is sent
- **THEN** input textarea is cleared

### Requirement: AI Response Display
The system SHALL display AI responses in the chat area as they stream from the server.

#### Scenario: Display AI response
- **WHEN** AIAnalysis message is received
- **THEN** analysis text appears in chat area aligned left

#### Scenario: Style AI messages
- **WHEN** AI message is rendered
- **THEN** it uses left-aligned styling distinct from user messages

### Requirement: Typing Indicator
The system SHALL show a typing indicator while waiting for AI response to generate.

#### Scenario: Show typing indicator
- **WHEN** question is sent and waiting for response
- **THEN** "AI is thinking..." with animated dots appears

#### Scenario: Hide typing indicator
- **WHEN** AI response starts streaming
- **THEN** typing indicator disappears and is replaced with response

### Requirement: Chat Message History
The system SHALL maintain scrollable history of chat messages.

#### Scenario: Store chat messages
- **WHEN** messages are exchanged
- **THEN** they are stored in model state

#### Scenario: Display message history
- **WHEN** TUI renders
- **THEN** recent messages are visible with older messages accessible via scroll
