layout (triangles) in;
layout (line_strip, max_vertices = 12) out;

// Model uniforms
uniform mat4 MVP;

// Inputs from Vertex Shader
in vec3 Normal[];

//// Inputs uniforms
//uniform int ShowWireframe;
//uniform int ShowVnormal;
//uniform int ShowFnormal;
// Const uniforms
const int ShowWireframe = 1;
const int ShowVnormal = 1;
const int ShowFnormal = 0;

// Colors
const vec4 colorWire    = vec4(1, 1, 0, 1);
const vec4 colorVnormal = vec4(1, 0, 0, 1);
const vec4 colorFnormal = vec4(0, 0, 1, 1);

// Output color to fragment shader
out vec4 Color;

void main() {

	// Emits triangle's vertices as lines to show wireframe
	if (ShowWireframe != 0) {
		for (int n = 0; n < gl_in.length(); n++) {
			// Vertex position
			gl_Position = MVP * gl_in[n].gl_Position;
			EmitVertex();
		}
		// Emit first triangle vertex to close the last line strip.
		gl_Position = MVP * gl_in[0].gl_Position;
		EmitVertex();
		EndPrimitive();
	}

	// Emits lines representing the vertices normals
	if (ShowVnormal != 0) {
		for (int i = 0; i < gl_in.length(); i++) {

			vec3 position = gl_in[i].gl_Position.xyz;
			vec3 normal = Normal[i];
			
			gl_Position = MVP * vec4(position, 1.0);
			Color = colorVnormal;
			EmitVertex();
			
			gl_Position = MVP * vec4(position + normal * 0.5, 1.0);
			Color = colorVnormal;
			EmitVertex();
			
			EndPrimitive();
		}
	}

	// Emits one line representing the face normal
	if (ShowFnormal != 0) {
		vec3 p0 = gl_in[0].gl_Position.xyz;
		vec3 p1 = gl_in[1].gl_Position.xyz;
		vec3 p2 = gl_in[2].gl_Position.xyz;
	  
		vec3 v0 = p0 - p1;
		vec3 v1 = p2 - p1;
		vec3 faceN = normalize(cross(v1, v0));

		// Center of the triangle
		vec3 center = (p0 + p1 + p2) / 3.0;
	  
		gl_Position = MVP * vec4(center, 1.0);
		Color = colorFnormal;
		EmitVertex();
	  
		gl_Position = MVP * vec4(center + faceN * 0.5, 1.0);
		Color = colorFnormal;
		EmitVertex();
		EndPrimitive();
	}
}
