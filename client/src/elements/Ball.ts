import { randomInt } from "../utils/random";

interface BallOption<T> {
  radius: number;
  canvasWidth: number;
  canvasHeight: number;
  content: T;
  createdAt: number;
}

class Ball<T = unknown> {
  public x: number;
  public y: number;
  public dx: number;
  public dy: number;
  public rotation: number;
  public opacity: number;
  public readonly radius: number;
  public readonly content: T;
  public readonly createdAt: number;

  constructor({ radius, canvasWidth, content, createdAt }: BallOption<T>) {
    this.x = randomInt(radius, canvasWidth - radius);
    this.y = -100;
    this.dx = randomInt(-2, 2);
    this.dy = randomInt(-2, 2);
    this.rotation = randomInt(0, 360);
    this.opacity = 1.0;
    this.radius = radius;
    this.content = content;
    this.createdAt = createdAt;
  }
}

export default Ball;
