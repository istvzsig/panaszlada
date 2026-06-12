export function isInHungary(lng, lat, bounds) {
  return (
    lng >= bounds[0][0] &&
    lng <= bounds[1][0] &&
    lat >= bounds[0][1] &&
    lat <= bounds[1][1]
  );
}
