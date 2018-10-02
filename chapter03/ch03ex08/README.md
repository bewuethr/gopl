`big.Rat` is so abismally slow that the output picture size is set to tiny.

For `complex64`, `complex128` and `big.Float`, the `runmb` script generates four pictures each at zoom levels around when artifacts start to occur. For the two complex types, this just means pixelation; `big.Float` suddenly creates an all blue output, not quite clear why.

And `big.Rat` is extremely slow at deeper zoom levels, so I haven't really tested its limits. Most likely, given that it is arbitrarily precise, the input using `bc` would become a problem for resolution before `big.Rat`.
