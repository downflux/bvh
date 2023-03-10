# bvh
BVH for the DownFlux server

This BVH struct is used for DownFlux collision detection. Objects stored in this
BVH must also specify the collision layers it will occupy. The concept of
collision layers is modeled after various game engines. The layer namespace is
user-defined. See
[pkg.go.dev/github.com/downflux/go-bvh](https://pkg.go.dev/github.com/downflux/go-bvh/bvh)
for the more general-purpose BVH data struct.

## Further Reading

* [Godot documentation](https://docs.godotengine.org/en/stable/tutorials/physics/physics_introduction.html)
