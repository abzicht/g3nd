
layout(local_size_x = 64) in; // Each workgroup processes 64 elements

// Define the buffer with multiple vectors
layout(std430, binding = 0) buffer DataBuffer {
    vec4 data[]; // Array of vec4s
};

void main() {
    uint index = gl_GlobalInvocationID.x; // Get the index of the current thread
    
    // Ensure the index is within bounds
    if (index < data.length()) {
        data[index] *= vec4(2.0, 0.5, 1.0, 1.0); // Example: Scale vector components
    }
}
