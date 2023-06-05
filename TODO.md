# TODO

- express computations (eg. mathematical, search, generation) in triples
    - [x] express mathematical unary function as a node, use it as a predicate with the argument is the object
    - [ ] express mathematical binary functions as reducers
        - eg: (x sum 1, x sum 2, sum = `(acc, x)->acc+x` ) -> (x equal 3)
        - sum is a reducer (accumulator, currentValue)

    - [ ] express queries as

    - [ ] generalize TripleSet.Compute() into functions that accept Triples

- setup a set for all nodes using the name of the class as a string prefix

- create a way to organize statements using concepts, words, ideas, dimensions, values about myself, my world

    - need a way to validate/invalidate statements (not every statement/question has meaning ... what is meaning?  statements that are promote life)
        - statements are
            - (subject predicate) pair
                - need a pair
                - need a way to validate a pair
                    - need a way to express what a valid pair is
                        - a valid pair is valid if pair.subject is in valid.subjects and ...
                            - a collection of valid pairs is defined by a pair of aggregates
        - need identites, properties, relationships
            - need aggregates, pairs
                - aggregates
                    - an item is either in or out
                        - need a way to express truth
                            - a property that promotes life (probably a circular definition, which is fine)
                            - an aggregate of true things?

                - pairs
                    - an item is either first, or second, or neither
                    - more generally, need ordered relationship
                        - an item is either before, or after, or neither
                            - need a way to express unknown
                                - an empty aggregate?
                                - or we define aggregates as: either in, out, or unknown
                                    - AggregateXXX

        - need a way to express boolean truth
            - organisms derive "truth" from similarity (similar is true, dissimilar is false)
                - organims derive "similar" from whether lived experience matches predicted experience (current response matches learned response)
                    - "does this item now consumed (therefore identified as resource) feels good or bad"
                        - organism learn through natural selection to respond to resources with "positive signal"
                        or
                        - organism increases metabolism after consuming a resource and that increased metabolism is interpreted as a "positive signal" (until it's no longer)


- create a way to generate all combinations of values given some dimension/category (eg. build a table header for color,size,style)


- relationships
    - similarity

    - resource/threat

    - aggregate
        - eg: (marc, alive)
        - set
        - pair

    - property
        - aggregate of different identities, indicating that they share the same property
		- with an expectation that the set is infinite

    - identity
        - aggregate of different properties, indicating that they refer to the same identity
        - with an expectation that this set is unique

