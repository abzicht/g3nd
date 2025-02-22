layout(local_size_x = 8, local_size_y = 8, local_size_z = 8) in;

layout(std430, shared, binding = 0) buffer DataBuffer1 {
    vec3 data[];
};
layout(std430, shared, binding = 1) buffer DataBuffer2 {
    uint i;
};
layout(std430, shared, binding = 2) buffer DataBuffer3 {
    vec3 loc;
    uint i;
    float speed;
    dvec4 data[];
} B;

void main() {
    uint index = gl_GlobalInvocationID.x; // Get the index of the current threa


    if (gl_GlobalInvocationID.yz == vec2(0)) {
        i+=10;
        B.i+=10;
        B.loc.xyz += vec3(1,2,3);
        B.speed -= 0.1;
        if (index < B.data.length()) {
            B.data[index].xyzw += dvec4(0.1,0.2,0.3,0.4);
        }
        if (index < data.length()) {
            data[index].xyz += vec3(0.1,0.2,0.3);
        }
    }
}
