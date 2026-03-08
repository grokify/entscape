/**
 * TypeScript types matching the entscape JSON schema.
 */

/** Package metadata */
export interface Package {
  name: string;
  source?: string;
  branch?: string;
  docs?: string;
}

/** Field attribute constants */
export type FieldAttr =
  | 'primary'
  | 'unique'
  | 'required'
  | 'optional'
  | 'immutable'
  | 'sensitive'
  | 'nillable'
  | 'default';

/** Entity field definition */
export interface Field {
  name: string;
  type: string;
  attrs?: FieldAttr[];
  default?: string;
  comment?: string;
}

/** Relation type constants */
export type RelationType = 'O2O' | 'O2M' | 'M2O' | 'M2M';

/** Entity edge (relationship) definition */
export interface Edge {
  name: string;
  target: string;
  relation: RelationType;
  inverse?: string;
  required?: boolean;
  unique?: boolean;
  comment?: string;
}

/** Database index definition */
export interface Index {
  fields: string[];
  unique?: boolean;
}

/** Entity definition */
export interface Entity {
  name: string;
  path?: string;
  description?: string;
  fields?: Field[];
  edges?: Edge[];
  indexes?: Index[];
  mixins?: string[];
}

/** Root schema structure */
export interface Schema {
  version: string;
  package?: Package;
  entities: Entity[];
}

/** Render options */
export interface RenderOptions {
  /** Layout direction: 'TB' (top-bottom), 'LR' (left-right) */
  direction?: 'TB' | 'LR';
  /** Enable source link click handlers */
  sourceLinks?: boolean;
  /** Custom node colors */
  colors?: {
    node?: string;
    nodeText?: string;
    nodeBorder?: string;
    edge?: string;
    edgeText?: string;
  };
  /** Callback when entity is clicked */
  onEntityClick?: (entity: Entity) => void;
  /** Callback when edge is clicked */
  onEdgeClick?: (edge: Edge, source: Entity, target: Entity) => void;
}

/** Entscape instance returned by render() */
export interface EntscapeInstance {
  /** Cytoscape core instance */
  cy: cytoscape.Core;
  /** Update with new schema data */
  update: (schema: Schema) => void;
  /** Fit view to all elements */
  fit: () => void;
  /** Zoom to specific entity */
  zoomTo: (entityName: string) => void;
  /** Highlight entity and its connections */
  highlight: (entityName: string) => void;
  /** Clear all highlights */
  clearHighlight: () => void;
  /** Destroy the instance */
  destroy: () => void;
}
