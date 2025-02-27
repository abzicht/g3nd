precision highp float;

// Buffers shared with a compute shader
layout(std430, shared, binding = 0) buffer DataBuffer1 {
    vec3 data[]; 
};
layout(std430, shared, binding = 1) buffer DataBuffer2 {
    uint i; 
};

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;
in vec2 VPosition;      // Vertex position in model coordinates (xy)

#include <lights>
#include <material>
#include <phong_model>

// Uniforms for configuration
uniform vec3 Color;

// Final fragment color
out vec4 FragColor;

vec3 rgb2hsv(vec3 c)
{
    vec4 K = vec4(0.0, -1.0 / 3.0, 2.0 / 3.0, -1.0);
    vec4 p = mix(vec4(c.bg, K.wz), vec4(c.gb, K.xy), step(c.b, c.g));
    vec4 q = mix(vec4(p.xyw, c.r), vec4(c.r, p.yzx), step(p.x, c.r));

    float d = q.x - min(q.w, q.y);
    float e = 1.0e-10;
    return vec3(abs(q.z + (q.w - q.y) / (6.0 * d + e)), d / (q.x + e), q.x);
}

// All components are in the range [0…1], including hue.
vec3 hsv2rgb(vec3 c)
{
    vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
    vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
    return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void main() {
    // Combine material with buffer from a compute shader
    vec2 index = VPosition;
    vec3 color = mix(Color, vec3(0), float(i%256)/255);
    color = mix(Color, hsv2rgb(vec3(0,1,1)), index.x * index.y);
    vec4 matDiffuse = vec4(color, 1.0);
    vec4 matAmbient = vec4(MatAmbientColor, MatOpacity) * vec4(color, 1.0);

    // Inverts the fragment normal if not FrontFacing
    vec3 fragNormal = Normal;
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, CamDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}
