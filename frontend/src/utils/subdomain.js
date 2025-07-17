export function getSubdomain() {
    const host = window.location.hostname; // e.g., "owner.localhost" or "www.owner.localhost"

    const parts = host.split('.'); // Split the hostname by '.'
  
    // Check if there are at least 2 parts (domain and TLD)
    if (parts.length >= 2) {
        // If the first part is 'www', return the second part as the subdomain
        if (parts[0] === 'www') {
            return parts.length > 2 ? parts[1] : null; // Return the second part if there's a subdomain
        }
        return parts[0]; // Return the first part as the subdomain if it's not 'www'
    } else {
        return null; // Return null if no subdomain
    }
}

export function getBaseDomain() {
    const host = window.location.hostname;

    // Handle localhost specifically
    if (host === 'localhost') {
        return 'localhost';
    }

    const parts = host.split('.');
    const numParts = parts.length;

    // For '4commander.my.id', 'my.id' is the TLD, and '4commander' is the domain.
    // So, we want '4commander.my.id'.
    // This handles cases like 'api.4commander.my.id', 'admin.4commander.my.id', or '4commander.my.id'
    if (numParts >= 3 && parts[numParts - 2] === 'my' && parts[numParts - 1] === 'id') {
        return parts.slice(numParts - 3).join('.'); // Returns '4commander.my.id'
    }

    // Default to the last two parts for standard domains (e.g., example.com from sub.example.com)
    return parts.slice(numParts - 2).join('.');
}