import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";

export const getUnitPrice = (unit: Unit) => unit.price;
