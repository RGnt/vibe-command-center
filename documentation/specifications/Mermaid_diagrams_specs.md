## Sequence diagrams

  |>  A Sequence diagram is and interaction diagram that shows how process operate with one another and in what order.

### Syntax

#### Participants

The participants can be defined implicitly. The participants or actors are rendered in order of appearance in the diagram soruce text. Sometimes you ming want the participants in a different order than how they appear in the first message. It is possible to specify actor's order of appearnce by doing the following:

```Mermaid
sequenceDiagram
    participant Alice
    participant Bob
    Bob->>Alice: Hi Alice
    Alice->>Bob: Hi bob
```

#### Actors 

If you specifically want to use the actor symbol instead of a rectange with text you can do so by using actor statements as per below:

```Mermaid
sequenceDiagram
    actor Alice
    actor Bob
    Alice->>Bob: Hi Bob
    Bob->>Alice: Hi Alice
```


#### Notes