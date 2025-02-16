
layout(local_size_x = 16) in; // Each workgroup processes 64 elements

// Define the buffer with multiple vectors
layout(std430, binding = 0) buffer DataBuffer {
    dvec4 data[]; // Array of vec4s
};
layout(std430, binding = 1) buffer DataBuffer2 {
    int is[];
};
layout(std430, binding = 2) buffer DataBuffer3 {
    int b[];
};

void main() {
    uint index = gl_GlobalInvocationID.x; // Get the index of the current thread
    
    // Ensure the index is within bounds
    if (index < data.length()) {
        data[index].xyzw = vec4(float(is[index]), data[index].x + data[index].y, data[index].z, data[index].w);
        //data[index].x = data[index].x + 1;
    }
    if( index <= is.length()) {
        is[index] += 1;
    }
    if( index < b.length()) {
        b[index+1] = b[index];
    }
}
