import jwt from "jsonwebtoken";

export interface Claims {
  userId: string;
  role: string;
}

export function verifyToken(token: string, secret: string): Claims {
  const decoded = jwt.verify(token, secret) as jwt.JwtPayload;

  if (typeof decoded.user_id !== "string" || decoded.user_id === "") {
    throw new Error("token invalido: falta user_id");
  }

  return {
    userId: decoded.user_id,
    role: typeof decoded.role === "string" ? decoded.role : "",
  };
}
