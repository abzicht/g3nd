#pragma optimize(off)
#pragma debug(on)
#pragma kernel curved
#pragma kernel straight

layout(local_size_x = 512, local_size_y = 1, local_size_z = 1) in;

layout(std430, binding = 0) buffer ParticlePos {
    vec3 positions[];
};
layout(std430, binding = 1) buffer ParticleColor {
    vec4 colors[];
};
//binding 2 is taken by particle shape
layout(std430, binding = 3) buffer ParticleVel {
    vec3 velocities[];
};

uniform vec3 BoundsMin;
uniform vec3 BoundsMax;

vec3 move_particle(vec3 position, vec3 velocity) {
    return position + velocity;
}
void set_color(uint id) {
    colors[id] = vec4(positions[id], 1);
}

vec3 bounce_on_bounds(uint id) {
    vec3 velocity = velocities[id];
    vec3 would_move = move_particle(positions[id], velocities[id]);
    if (would_move.x < BoundsMin.x  || would_move.x > BoundsMax.x) {
        velocity.x *= -1;
    }
    if (would_move.y < BoundsMin.z  || would_move.y > BoundsMax.y) {
        velocity.y *= -1;
    }
    if (would_move.z < BoundsMin.z  || would_move.z > BoundsMax.z) {
        velocity.z *= -1;
    }
    return velocity;
}

void straight() {
    uint id = gl_GlobalInvocationID.x;
    if (id > positions.length()) return; // Safety check
    if (id > velocities.length()) return; // Safety check

    velocities[id] = bounce_on_bounds(id);
    positions[id] = move_particle(positions[id], velocities[id]);
    set_color(id);
}

void curved() {
    uint id = gl_GlobalInvocationID.x;
    if (id > positions.length()) return; // Safety check
    if (id > velocities.length()) return; // Safety check

    velocities[id] = bounce_on_bounds(id);
    velocities[id].y -= 0.001;
    positions[id] = move_particle(positions[id], velocities[id]);
    set_color(id);
}
