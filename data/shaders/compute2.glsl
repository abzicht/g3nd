layout(local_size_x = 16) in; // Each workgroup processes 16 elements

// See https://www.khronos.org/opengl/wiki/Interface_Block_(GLSL) for
// buffer layouts and their configuration

// Define the buffer with multiple vectors
layout(std430, shared, binding = 1) buffer DataBuffer1 {
    dvec4 data[]; // Array of vec4s
};
layout(std430, shared, binding = 2) buffer DataBuffer2 {
    int is[];
};
layout(std430, shared, binding = 3) buffer DataBuffer3 {
    int b[];
};

struct T {
    bool isactive;
    float location_;
    uint speed;
};

layout(std430, shared, binding = 4) buffer MyTees {
    T tees[];
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
    if( index < tees.length()) {
        tees[index].isactive = true;
        tees[index].speed += 1;
        tees[index].location_ = 1;
    }
}
