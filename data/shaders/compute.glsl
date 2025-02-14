layout(std430, binding = 0) buffer DataBuffer {
    vec4 result;
};

layout (local_size_x = 4, local_size_y = 1, local_size_z = 1) in;

void main() {
    result = vec4(1.0, 2.0, 3.0, 4.0); // Example output data
}
