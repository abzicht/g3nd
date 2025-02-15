
layout(local_size_x = 16) in; // Each workgroup processes 64 elements

// Define the buffer with multiple vectors
layout(std430, binding = 0) buffer DataBuffer {
    vec4 data[]; // Array of vec4s
};
layout(std430, binding = 1) buffer DataBuffer2 {
    int i;
};

void main() {
    uint index = gl_GlobalInvocationID.x; // Get the index of the current thread
    
    // Ensure the index is within bounds
    if (index < data.length()) {
        data[index].xyzw = vec4(0, data[index].x + data[index].y, data[index].z, data[index].w);
        //data[index].xyzw = vec4(data[index].x * -1, 0, 0, 0);
    }
    if( index == 0) {
        i += 1;
    }
}
