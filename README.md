# GPN Client

A client implementation for the tron game that was playable at GPN 2023: https://github.com/freehuntx/gpn-tron

I only participated remotely back then and my original implementation was lost. However, since we did a small coding excercise with
this project at work as well, I decided to try and rewrite it again.

## General structure

The main loop and general plumbing of the network protocol is happening in `main.go`.

The interesting file to implement a new behaviour is `decider.go`, which requires you to implement the `Decider` interface for
a given bot. The `Decide` method receives the current world state as an argument and has to output the next movement direction
in response.

`world.go` implements all things representing current world state and some basic helpers to "see" the world, such as finding the
position of a cell in a given direction and finding neighbour cells.
