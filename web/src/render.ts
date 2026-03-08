import cytoscape, { Core, ElementDefinition, NodeDataDefinition, EdgeDataDefinition } from 'cytoscape';
import cytoscapeDagre from 'cytoscape-dagre';
import type { Schema, Entity, Edge, RenderOptions, EntscapeInstance } from './types';
import { getStylesheet, relationLabels } from './styles';

// Register dagre layout
let dagreRegistered = false;
function registerDagre() {
  if (!dagreRegistered && typeof cytoscape !== 'undefined') {
    cytoscape.use(cytoscapeDagre);
    dagreRegistered = true;
  }
}

/** Format field for display */
function formatField(field: { name: string; type: string; attrs?: string[] }): string {
  const attrs = field.attrs || [];
  const markers: string[] = [];

  if (attrs.includes('primary')) markers.push('PK');
  if (attrs.includes('unique')) markers.push('U');
  if (attrs.includes('required')) markers.push('*');

  const markerStr = markers.length > 0 ? ` [${markers.join(',')}]` : '';
  return `${field.name}: ${field.type}${markerStr}`;
}

/** Build node content (fields list) */
function buildNodeContent(entity: Entity): string {
  if (!entity.fields || entity.fields.length === 0) {
    return '';
  }
  return entity.fields.map(formatField).join('\n');
}

/** Calculate node height based on field count */
function calculateNodeHeight(entity: Entity): number {
  const baseHeight = 50;
  const fieldHeight = 18;
  const fieldCount = entity.fields?.length || 0;
  return baseHeight + fieldCount * fieldHeight;
}

/** Convert schema to Cytoscape elements */
function schemaToElements(schema: Schema): ElementDefinition[] {
  const elements: ElementDefinition[] = [];
  const entityMap = new Map<string, Entity>();

  // Build entity lookup
  for (const entity of schema.entities) {
    entityMap.set(entity.name, entity);
  }

  // Create nodes
  for (const entity of schema.entities) {
    const nodeData: NodeDataDefinition = {
      id: entity.name,
      label: entity.name,
      content: buildNodeContent(entity),
      height: calculateNodeHeight(entity),
      entity: entity,
    };

    if (entity.path) {
      nodeData.path = entity.path;
    }

    elements.push({ data: nodeData });
  }

  // Create edges
  const edgeSet = new Set<string>();

  for (const entity of schema.entities) {
    if (!entity.edges) continue;

    for (const edge of entity.edges) {
      // Skip if target entity doesn't exist
      if (!entityMap.has(edge.target)) continue;

      // Create unique edge ID to avoid duplicates
      const edgeId = `${entity.name}-${edge.name}-${edge.target}`;
      const reverseId = `${edge.target}-${edge.inverse || ''}-${entity.name}`;

      // Skip if we've already added this edge (for inverse relationships)
      if (edgeSet.has(reverseId)) continue;
      edgeSet.add(edgeId);

      const edgeData: EdgeDataDefinition = {
        id: edgeId,
        source: entity.name,
        target: edge.target,
        label: relationLabels[edge.relation] || edge.relation,
        name: edge.name,
        relation: edge.relation,
        required: edge.required || false,
        optional: !edge.required,
        edge: edge,
        sourceEntity: entity,
        targetEntity: entityMap.get(edge.target),
      };

      elements.push({ data: edgeData });
    }
  }

  return elements;
}

/** Create entscape instance */
export function render(
  container: HTMLElement | string,
  schema: Schema,
  options: RenderOptions = {}
): EntscapeInstance {
  const containerEl =
    typeof container === 'string' ? document.querySelector(container) : container;

  if (!containerEl) {
    throw new Error(`Container not found: ${container}`);
  }

  // Ensure dagre is registered
  registerDagre();

  const direction = options.direction || 'TB';
  const elements = schemaToElements(schema);

  // Create Cytoscape instance
  const cy: Core = cytoscape({
    container: containerEl as HTMLElement,
    elements,
    style: getStylesheet(options.colors),
    layout: {
      name: 'dagre',
      rankDir: direction,
      nodeSep: 80,
      rankSep: 100,
      edgeSep: 50,
      padding: 30,
    } as cytoscape.LayoutOptions,
    minZoom: 0.1,
    maxZoom: 3,
    wheelSensitivity: 0.3,
  });

  // Set up click handlers
  if (options.sourceLinks !== false) {
    cy.on('tap', 'node', (evt) => {
      const node = evt.target;
      const entity = node.data('entity') as Entity;

      if (options.onEntityClick) {
        options.onEntityClick(entity);
      } else if (entity.path) {
        window.open(entity.path, '_blank');
      }
    });
  }

  if (options.onEdgeClick) {
    const onEdgeClick = options.onEdgeClick;
    cy.on('tap', 'edge', (evt) => {
      const edgeEl = evt.target;
      const edge = edgeEl.data('edge') as Edge;
      const sourceEntity = edgeEl.data('sourceEntity') as Entity;
      const targetEntity = edgeEl.data('targetEntity') as Entity;
      onEdgeClick(edge, sourceEntity, targetEntity);
    });
  }

  // Instance methods
  const instance: EntscapeInstance = {
    cy,

    update(newSchema: Schema) {
      const newElements = schemaToElements(newSchema);
      cy.elements().remove();
      cy.add(newElements);
      cy.layout({
        name: 'dagre',
        rankDir: direction,
        nodeSep: 80,
        rankSep: 100,
        edgeSep: 50,
        padding: 30,
      } as cytoscape.LayoutOptions).run();
    },

    fit() {
      cy.fit(undefined, 30);
    },

    zoomTo(entityName: string) {
      const node = cy.getElementById(entityName);
      if (node.length > 0) {
        cy.animate({
          center: { eles: node },
          zoom: 1.5,
          duration: 300,
        });
      }
    },

    highlight(entityName: string) {
      instance.clearHighlight();
      const node = cy.getElementById(entityName);
      if (node.length > 0) {
        node.addClass('highlighted');
        node.connectedEdges().addClass('highlighted');
        node.neighborhood('node').addClass('highlighted');
      }
    },

    clearHighlight() {
      cy.elements().removeClass('highlighted');
    },

    destroy() {
      cy.destroy();
    },
  };

  return instance;
}
