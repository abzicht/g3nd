precision highp float;

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;
in vec2 VPosition;      // Vertex position in model coordinates (xy)

#include <lights>
#include <material>
#include <phong_model>

// Uniforms for configure brick pattern
uniform vec3 Color;
uniform vec2 Size;

// Final fragment color
out vec4 FragColor;

void main() {

    vec2 position = VPosition;
    vec3 color = mix(Color, vec3(0,0,0), position.x*position.y);

    // Combine material with brick pattern colors
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

