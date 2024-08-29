# veil
Veil simplifies RPC code while trying to not be "magical".  It consist of two components.

* a code generator that uses templates to generate bindings between a Go struct and a RPC library.  By default it binds to net/RPC but I have plans to support others in the future.  
* a library used to register structs to be exposed via RPC and calls to lookup client implementations that make calls using that RPC to the exposed service.

