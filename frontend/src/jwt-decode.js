// Pequeno utilitário para decodificar JWT (sem dependência externa)
export default function jwt_decode(token) {
  if (!token) return {};
  try {
    const payload = token.split('.')[1];
    const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
    return decoded;
  } catch {
    return {};
  }
}
