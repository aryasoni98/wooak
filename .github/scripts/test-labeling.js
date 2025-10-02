#!/usr/bin/env node

/**
 * Test script for auto-labeling functionality
 * This script simulates the label extraction logic
 */

// Test data - simulate issue body content
const testIssueBodies = [
  {
    name: "Bug Fix Issue",
    body: `🎉 Welcome to Hacktoberfest! Thank you for considering contributing to Wooak.

Contribution Type: Bug Fix
Difficulty Level: Beginner - Good first issue

Contribution Description: Fix the SSH connection timeout issue

Motivation: I want to help improve the SSH management tool

Relevant Experience: I have experience with Go and SSH protocols`,
    expectedLabels: ['hacktoberfest', 'hacktoberfest-accepted', 'bug']
  },
  {
    name: "New Feature Issue", 
    body: `🎉 Welcome to Hacktoberfest! Thank you for considering contributing to Wooak.

Contribution Type: New Feature
Difficulty Level: Intermediate - Some experience needed

Contribution Description: Add dark mode support to the TUI

Motivation: Dark mode would improve user experience

Relevant Experience: I have experience with Go and TUI libraries`,
    expectedLabels: ['hacktoberfest', 'hacktoberfest-accepted', 'enhancement']
  },
  {
    name: "Documentation Issue",
    body: `🎉 Welcome to Hacktoberfest! Thank you for considering contributing to Wooak.

Contribution Type: Documentation
Difficulty Level: Beginner - Good first issue

Contribution Description: Improve the README with better installation instructions

Motivation: I want to help new users get started easily

Relevant Experience: I have experience with technical writing`,
    expectedLabels: ['hacktoberfest', 'hacktoberfest-accepted', 'documentation']
  }
];

// Helper function to extract contribution type from issue body
function extractContributionType(body) {
  if (!body) return null;
  
  // Look for the contribution type dropdown value
  const contributionMatch = body.match(/Contribution Type[:\s]*([^\n\r]+)/i);
  if (contributionMatch) {
    return contributionMatch[1].trim();
  }
  
  return null;
}

// Helper function to map contribution type to label
function getContributionTypeLabel(contributionType) {
  const typeMapping = {
    'Bug Fix': 'bug',
    'New Feature': 'enhancement',
    'Documentation': 'documentation',
    'Code Refactoring': 'refactor',
    'Test Coverage': 'test',
    'Performance Improvement': 'performance',
    'Security Enhancement': 'security',
    'UI/UX Improvement': 'ui/ux'
  };
  
  return typeMapping[contributionType] || null;
}

// Test the extraction logic
function testLabeling() {
  console.log('🧪 Testing Auto-labeling Logic\n');
  
  testIssueBodies.forEach((testCase, index) => {
    console.log(`Test ${index + 1}: ${testCase.name}`);
    console.log('─'.repeat(50));
    
    const contributionType = extractContributionType(testCase.body);
    const typeLabel = getContributionTypeLabel(contributionType);
    
    const actualLabels = ['hacktoberfest', 'hacktoberfest-accepted'];
    if (typeLabel) {
      actualLabels.push(typeLabel);
    }
    
    console.log(`Contribution Type: ${contributionType}`);
    console.log(`Type Label: ${typeLabel}`);
    console.log(`Expected Labels: [${testCase.expectedLabels.join(', ')}]`);
    console.log(`Actual Labels: [${actualLabels.join(', ')}]`);
    
    const isCorrect = JSON.stringify(actualLabels.sort()) === JSON.stringify(testCase.expectedLabels.sort());
    console.log(`Result: ${isCorrect ? '✅ PASS' : '❌ FAIL'}`);
    console.log('');
  });
}

// Run the test
testLabeling();
