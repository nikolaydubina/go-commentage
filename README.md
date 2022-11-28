# go-commentage

> Are my comments being updated together with code?

Inspired by "Clean Code" by Robert C. Martin, this tool collects details on how far comments drift from code they describe.

### Requirements

You need to have `git` with version `>=2.37`.

```bash
go install github.com/nikolaydubina/go-commentage
```

### Examples

#### [kubernetes](https://github.com/kubernetes/kubernetes)

TODO

### Heuristics

#### Simple Age Difference

Measure of how far away in terms of days or commits last update of function body as compared to last update to associated doc comment.

#### Weighted Age Difference

> Work in Progress

Code changes happen at various rates.
Comments and code can change one line or can change 90% of lines.
It is useful to differentiate between updates.

### References

* https://git-scm.com/docs/git-blame
* https://github.com/nishanths/exhaustive
* https://github.com/kubernetes/kubernetes
