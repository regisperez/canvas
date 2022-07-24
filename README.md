# ASCII canvas

Your task is to build a client-server system to represent an ASCII art drawing canvas. The exercise involves two pieces:

- A server that will implement the drawing endpoints into a canvas.
- A _read-only_ client that will show the canvas stored by the server with the drawings.

## What we are looking for is...

- correctly functioning solution that comes with running instructions and examples.
- readable, clear code of the kind you would yourself like your team-mates to write and maintain with you.
- appropriately automatically tested, including failure scenarios.
- appropriate data structures and algorithms being applied as part of your solution.
- a version history being included with your solution.

### The server should be...

- implementing the canvas behaviours described below into a web service API as you see fit.
- written in Elixir if you are applying to an Elixir position, and in Go, if you are applying to a Go position. In case this a a mixed position, please, clarify with the recruiter.
- automatically tested appropriately.
- runnable with clear running instructions.

### Canvas should be...

- identifiable with a globally unique identifier.
- persisted across application launches.
- **_not_** authenticating users: authenticating or authorizing requests is out of scope.

### Drawing operation

- A rectangle parameterised with...
  - Coordinates for the **upper-left corner**.
  - **width** and **height**.
  - an optional **fill** character.
  - an optional **outline** character.
  - One of either **fill** or **outline** should always be present.

A character can be assumed to be an ASCII encoded byte.

The canvas can be assumed to be a fixed size.

### The read-only client should be...

- Non-interactive: you do **_not_** need to implement any client-side interactions. Read-only.
- simple: you can implement the drawing operations with constant-width characters.
- runnable with clear running instructions.

### Test fixture 1

- Rectangle at [3,2] with width: 5, height: 3, outline character: `@`, fill character: `X`
- Rectangle at [10, 3] with width: 14, height: 6, outline character: `X`, fill character: `O`

```


   @@@@@
   @XXX@  XXXXXXXXXXXXXX
   @@@@@  XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XXXXXXXXXXXXXX
```

### Test fixture 2

- Rectangle at `[14, 0]` with width `7`, height `6`, outline character: none, fill: `.`
- Rectangle at `[0, 3]` with width `8`, height `4`, outline character: `O`, fill: `none`
- Rectangle at `[5, 5]` with width `5`, height `3`, outline character: `X`, fill: `X`

```
              .......
              .......
              .......
OOOOOOOO      .......
O      O      .......
O    XXXXX    .......
OOOOOXXXXX
     XXXXX
```
