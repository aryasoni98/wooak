#!/usr/bin/env node

/**
 * PR Title Validator
 * Validates pull request titles against conventional commit format
 */

const fs = require('fs');
const path = require('path');

// Configuration
const CONFIG = {
  types: [
    'feat', 'fix', 'docs', 'style', 'refactor', 'perf', 'test', 
    'chore', 'ci', 'build', 'upgrade', 'update', 'bump', 'security', 'deps'
  ],
  scopes: [
    'ui', 'cli', 'config', 'parser', 'deps', 'security', 
    'workflow', 'makefile', 'docker', 'ci'
  ],
  hacktoberfestPrefix: '[Hacktoberfest]:'
};

// Test cases
const testCases = [
  {
    title: 'feat: add dark mode support to TUI',
    expected: true,
    description: 'Basic feature with conventional format'
  },
  {
    title: 'fix(ui): resolve connection timeout issue',
    expected: true,
    description: 'Bug fix with scope'
  },
  {
    title: '[Hacktoberfest]: feat: add macOS support for Ollama',
    expected: true,
    description: 'Hacktoberfest PR with conventional format'
  },
  {
    title: 'upgrade: bump Go version to 1.21',
    expected: true,
    description: 'Upgrade type'
  },
  {
    title: 'Upgrade Codebase',
    expected: false,
    description: 'Missing type prefix (the failing case)'
  },
  {
    title: 'fix: resolve bug',
    expected: true,
    description: 'Simple bug fix'
  },
  {
    title: 'docs: update installation guide',
    expected: true,
    description: 'Documentation update'
  },
  {
    title: 'chore: update dependencies',
    expected: true,
    description: 'Maintenance task'
  },
  {
    title: 'security: fix vulnerability in key validation',
    expected: true,
    description: 'Security fix'
  },
  {
    title: 'perf: optimize server connection pooling',
    expected: true,
    description: 'Performance improvement'
  }
];

/**
 * Validates a PR title against conventional commit format
 * @param {string} title - The PR title to validate
 * @returns {object} - Validation result with success, errors, and suggestions
 */
function validatePRTitle(title) {
  const result = {
    success: false,
    errors: [],
    suggestions: [],
    type: null,
    scope: null,
    description: null,
    isHacktoberfest: false
  };

  if (!title || typeof title !== 'string') {
    result.errors.push('Title is required and must be a string');
    return result;
  }

  const trimmedTitle = title.trim();
  
  // Check for Hacktoberfest prefix
  if (trimmedTitle.startsWith(CONFIG.hacktoberfestPrefix)) {
    result.isHacktoberfest = true;
    const withoutPrefix = trimmedTitle.substring(CONFIG.hacktoberfestPrefix.length).trim();
    return validateConventionalFormat(withoutPrefix, result);
  }

  return validateConventionalFormat(trimmedTitle, result);
}

/**
 * Validates conventional commit format
 * @param {string} title - The title to validate
 * @param {object} result - The result object to update
 * @returns {object} - Updated result object
 */
function validateConventionalFormat(title, result) {
  // Pattern: type(scope): description
  const conventionalPattern = /^([a-z]+)(?:\(([a-z-]+)\))?:\s*(.+)$/;
  const match = title.match(conventionalPattern);

  if (!match) {
    result.errors.push('Title must follow conventional commit format: "type: description" or "type(scope): description"');
    result.suggestions.push('Examples: "feat: add new feature", "fix(ui): resolve bug", "docs: update guide"');
    return result;
  }

  const [, type, scope, description] = match;

  // Validate type
  if (!CONFIG.types.includes(type)) {
    result.errors.push(`Invalid type "${type}". Must be one of: ${CONFIG.types.join(', ')}`);
    result.suggestions.push(`Common types: feat, fix, docs, chore, upgrade`);
  }

  // Validate scope (if provided)
  if (scope && !CONFIG.scopes.includes(scope)) {
    result.errors.push(`Invalid scope "${scope}". Must be one of: ${CONFIG.scopes.join(', ')}`);
    result.suggestions.push(`Common scopes: ui, cli, config, deps`);
  }

  // Validate description
  if (!description || description.length < 3) {
    result.errors.push('Description must be at least 3 characters long');
  }

  if (description && description[0] === description[0].toUpperCase()) {
    result.errors.push('Description should start with lowercase letter');
    result.suggestions.push(`Change "${description}" to "${description[0].toLowerCase() + description.slice(1)}"`);
  }

  if (description && description.endsWith('.')) {
    result.errors.push('Description should not end with a period');
    result.suggestions.push(`Remove the period from "${description}"`);
  }

  // Set parsed values
  result.type = type;
  result.scope = scope;
  result.description = description;
  result.success = result.errors.length === 0;

  return result;
}

/**
 * Provides suggestions for fixing a title
 * @param {string} title - The original title
 * @param {object} validation - The validation result
 * @returns {string[]} - Array of suggested titles
 */
function getSuggestions(title, validation) {
  const suggestions = [];

  if (!validation.success) {
    // Try to fix common issues
    const lowerTitle = title.toLowerCase();
    
    // Check if it's a simple case of missing prefix
    if (lowerTitle.includes('upgrade') || lowerTitle.includes('update')) {
      suggestions.push(`upgrade: ${title.toLowerCase()}`);
    } else if (lowerTitle.includes('fix') || lowerTitle.includes('bug')) {
      suggestions.push(`fix: ${title.toLowerCase()}`);
    } else if (lowerTitle.includes('add') || lowerTitle.includes('new') || lowerTitle.includes('feature')) {
      suggestions.push(`feat: ${title.toLowerCase()}`);
    } else if (lowerTitle.includes('doc') || lowerTitle.includes('readme')) {
      suggestions.push(`docs: ${title.toLowerCase()}`);
    } else {
      suggestions.push(`chore: ${title.toLowerCase()}`);
    }

    // Add Hacktoberfest version if applicable
    if (validation.isHacktoberfest) {
      suggestions.push(`[Hacktoberfest]: ${suggestions[0]}`);
    }
  }

  return suggestions;
}

/**
 * Runs the test suite
 */
function runTests() {
  console.log('🧪 Testing PR Title Validator\n');
  
  let passed = 0;
  let failed = 0;

  testCases.forEach((testCase, index) => {
    const result = validatePRTitle(testCase.title);
    const success = result.success === testCase.expected;
    
    console.log(`Test ${index + 1}: ${testCase.description}`);
    console.log(`Title: "${testCase.title}"`);
    console.log(`Expected: ${testCase.expected ? 'Valid' : 'Invalid'}`);
    console.log(`Actual: ${result.success ? 'Valid' : 'Invalid'}`);
    
    if (!success) {
      console.log(`❌ FAIL`);
      if (result.errors.length > 0) {
        console.log(`Errors: ${result.errors.join(', ')}`);
      }
      if (result.suggestions.length > 0) {
        console.log(`Suggestions: ${result.suggestions.join(', ')}`);
      }
      failed++;
    } else {
      console.log(`✅ PASS`);
      passed++;
    }
    console.log('');
  });

  console.log(`\n📊 Test Results: ${passed} passed, ${failed} failed`);
  return failed === 0;
}

/**
 * Main function
 */
function main() {
  const args = process.argv.slice(2);
  
  if (args.length === 0) {
    // Run tests
    const success = runTests();
    process.exit(success ? 0 : 1);
  } else {
    // Validate provided title
    const title = args.join(' ');
    const result = validatePRTitle(title);
    
    console.log(`🔍 Validating PR Title: "${title}"\n`);
    
    if (result.success) {
      console.log('✅ Title is valid!');
      console.log(`Type: ${result.type}`);
      if (result.scope) console.log(`Scope: ${result.scope}`);
      console.log(`Description: ${result.description}`);
      if (result.isHacktoberfest) console.log('🎃 Hacktoberfest PR detected');
    } else {
      console.log('❌ Title is invalid!');
      console.log('\nErrors:');
      result.errors.forEach(error => console.log(`  - ${error}`));
      
      if (result.suggestions.length > 0) {
        console.log('\nSuggestions:');
        result.suggestions.forEach(suggestion => console.log(`  - ${suggestion}`));
      }
      
      const titleSuggestions = getSuggestions(title, result);
      if (titleSuggestions.length > 0) {
        console.log('\nSuggested titles:');
        titleSuggestions.forEach(suggestion => console.log(`  - "${suggestion}"`));
      }
    }
    
    process.exit(result.success ? 0 : 1);
  }
}

// Run if called directly
if (require.main === module) {
  main();
}

module.exports = {
  validatePRTitle,
  getSuggestions,
  runTests
};
