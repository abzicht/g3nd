layout(local_size_x = 8, local_size_y = 8, local_size_z = 8) in;

layout(std430, shared, binding = 0) buffer VertexPos {
    vec3 positions[];
};

uniform int gridWidth;
uniform int gridHeight;

void main() {
    uint x = gl_GlobalInvocationID.x;
    uint y = gl_GlobalInvocationID.y;
    uint z = gl_GlobalInvocationID.z;

    uint id = z * gridWidth * gridHeight + y * gridWidth + x;
    if (id >= positions.length()) return; // Safety check

    positions[id].xyz += vec3(0.001);
    //positions[index_].xyz = vec3(1);
}
