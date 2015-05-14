using Go = import "go.capnp";
$Go.package("driver");
$Go.import("driver");

@0xfbb40296665d9a47;

struct Z {
    jobVec  @0: List(Job);
}

struct Job {
    id      @0: Int64;
    name    @1: Text;
    func    @2: Text;
    args    @3: Text;
    timeout @4: Int64;
    schedAt @5: Int64;
    runAt   @6: Int64;
    status  @7: Text;
}
