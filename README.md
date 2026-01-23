# Zugzwang

**Zugzwang** is a self-led, senior project: a real-time online chess platform featuring a Go-based backend, a modern web frontend, and a custom-built chess engine.

You view the project here [zugzwang.dev](https://zugzwang.dev).

## Collaborators
- Liam Blanton
- Daniel Hansen

## Goals
The primary goals of this project are to:
- **Design, build, and deploy** a production-quality web application that demonstrates core computer science fundamentals
- Strengthen skills in **Go, TypeScript, React, and AWS**, with an emphasis on concurrency, networking, and state management
- Built a clean, modern frontend that users can play chess on 
- Implement a competitive chess engine with a target playing strength of **~2000 ELO**
- Track basic anaylitics about our users

## Requirements
### Functional
- Users can play a game of chess against our engine
- The system enforces all standard chess rules, including legal move validation and game-ending conditions
- The system maintains authoritative game state on the server
- The frontend displays the board state, legal moves, clocks, and game outcome
- The system records basic game and user interaction data for analytics purposes
### Non-Functional
- The system provides near real-time game updates between the player and the server
- Backend services safely handle concurrent game state updates
- The application remains responsive across desktop and mobile devices
- Core chess logic is isolated from networking and UI code for maintainability
- The system is deployable to a cloud environment (AWS)
- The codebase is structured to support future feature expansion


