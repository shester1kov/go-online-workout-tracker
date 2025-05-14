import { Role } from "../models/roles";

export const hasRole = (roles: Role[], roleName: string): boolean => {
  return roles.some(role => role.name === roleName);
};
