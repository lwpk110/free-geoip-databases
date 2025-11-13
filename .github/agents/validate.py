#!/usr/bin/env python3
"""
Validate custom agent configurations in .github/agents/

This script checks that all YAML files in the agents directory are:
1. Valid YAML syntax
2. Contain required fields
3. Follow the expected structure
"""

import sys
import os
from pathlib import Path

try:
    import yaml
except ImportError:
    print("Installing PyYAML...")
    import subprocess
    subprocess.check_call([sys.executable, "-m", "pip", "install", "-q", "pyyaml"])
    import yaml


def validate_agent_config(filepath):
    """Validate a single agent configuration file."""
    errors = []
    warnings = []
    
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            config = yaml.safe_load(f)
        
        # Check required fields
        required_fields = ['name', 'description']
        for field in required_fields:
            if field not in config:
                errors.append(f"Missing required field: {field}")
        
        # Check recommended fields
        recommended_fields = ['metadata', 'expertise', 'behavior', 'instructions']
        for field in recommended_fields:
            if field not in config:
                warnings.append(f"Missing recommended field: {field}")
        
        # Validate structure
        if 'name' in config and not isinstance(config['name'], str):
            errors.append("Field 'name' must be a string")
        
        if 'description' in config and not isinstance(config['description'], str):
            errors.append("Field 'description' must be a string")
        
        # Check metadata structure if present
        if 'metadata' in config:
            if not isinstance(config['metadata'], dict):
                errors.append("Field 'metadata' must be a dictionary")
        
        return {
            'valid': len(errors) == 0,
            'errors': errors,
            'warnings': warnings,
            'config': config
        }
        
    except yaml.YAMLError as e:
        return {
            'valid': False,
            'errors': [f"YAML parsing error: {str(e)}"],
            'warnings': [],
            'config': None
        }
    except Exception as e:
        return {
            'valid': False,
            'errors': [f"Unexpected error: {str(e)}"],
            'warnings': [],
            'config': None
        }


def main():
    """Main validation function."""
    agents_dir = Path('.github/agents')
    
    if not agents_dir.exists():
        print(f"❌ Agents directory not found: {agents_dir}")
        return 1
    
    # Find all YAML files
    yaml_files = list(agents_dir.glob('*.yml')) + list(agents_dir.glob('*.yaml'))
    
    if not yaml_files:
        print(f"⚠️  No agent configuration files found in {agents_dir}")
        return 0
    
    print(f"🔍 Validating {len(yaml_files)} agent configuration(s)...\n")
    
    all_valid = True
    
    for filepath in yaml_files:
        print(f"📄 Checking: {filepath.name}")
        result = validate_agent_config(filepath)
        
        if result['valid']:
            print(f"✅ Valid configuration")
            if result['config'] and 'name' in result['config']:
                print(f"   Agent Name: {result['config']['name']}")
        else:
            print(f"❌ Invalid configuration")
            all_valid = False
        
        if result['errors']:
            for error in result['errors']:
                print(f"   ❌ Error: {error}")
        
        if result['warnings']:
            for warning in result['warnings']:
                print(f"   ⚠️  Warning: {warning}")
        
        print()
    
    if all_valid:
        print("✅ All agent configurations are valid!")
        return 0
    else:
        print("❌ Some agent configurations have errors")
        return 1


if __name__ == '__main__':
    sys.exit(main())
