import type cytoscape from 'cytoscape';

/** Default color palette */
export const defaultColors = {
  node: '#f8f9fa',
  nodeText: '#212529',
  nodeBorder: '#dee2e6',
  nodeHeaderBg: '#e9ecef',
  edge: '#6c757d',
  edgeText: '#495057',
  highlight: '#0d6efd',
  highlightBg: '#cfe2ff',
};

/** Relation type to display label */
export const relationLabels: Record<string, string> = {
  O2O: '1:1',
  O2M: '1:N',
  M2O: 'N:1',
  M2M: 'M:N',
};

/** Generate Cytoscape stylesheet */
export function getStylesheet(colors: Partial<typeof defaultColors> = {}): cytoscape.StylesheetStyle[] {
  const c = { ...defaultColors, ...colors };

  return [
    // Entity nodes
    {
      selector: 'node',
      style: {
        'background-color': c.node,
        'border-color': c.nodeBorder,
        'border-width': 1,
        'border-style': 'solid',
        label: 'data(label)',
        'text-valign': 'top',
        'text-halign': 'center',
        'text-margin-y': 8,
        'font-family': 'system-ui, -apple-system, sans-serif',
        'font-size': 14,
        'font-weight': 'bold',
        color: c.nodeText,
        shape: 'round-rectangle',
        width: 'label',
        height: 'data(height)',
        'padding-top': '30px',
        'padding-bottom': '10px',
        'padding-left': '15px',
        'padding-right': '15px',
        'text-wrap': 'wrap',
        'text-max-width': '200px',
      },
    },
    // Highlighted nodes
    {
      selector: 'node.highlighted',
      style: {
        'background-color': c.highlightBg,
        'border-color': c.highlight,
        'border-width': 2,
      },
    },
    // Edges
    {
      selector: 'edge',
      style: {
        width: 2,
        'line-color': c.edge,
        'target-arrow-color': c.edge,
        'target-arrow-shape': 'triangle',
        'curve-style': 'bezier',
        label: 'data(label)',
        'font-family': 'system-ui, -apple-system, sans-serif',
        'font-size': 11,
        color: c.edgeText,
        'text-background-color': '#ffffff',
        'text-background-opacity': 0.9,
        'text-background-padding': '3px',
        'text-rotation': 'autorotate',
      },
    },
    // Highlighted edges
    {
      selector: 'edge.highlighted',
      style: {
        'line-color': c.highlight,
        'target-arrow-color': c.highlight,
        width: 3,
      },
    },
    // Required edges (solid)
    {
      selector: 'edge[required]',
      style: {
        'line-style': 'solid',
      },
    },
    // Optional edges (dashed)
    {
      selector: 'edge[?optional]',
      style: {
        'line-style': 'dashed',
      },
    },
    // M2M edges (double arrow)
    {
      selector: 'edge[relation = "M2M"]',
      style: {
        'source-arrow-shape': 'triangle',
        'source-arrow-color': c.edge,
      },
    },
  ];
}
