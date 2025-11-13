# GitHub Copilot Agents

This directory contains custom GitHub Copilot agent configurations for the free-geoip-databases project.

## Available Agents

### Senior Coding Engineer (资深编码工程师)

**Files:**
- `senior-coding-engineer.yml` (English)
- `senior-coding-engineer-zh.yml` (中文)

**Purpose:**
A highly experienced senior development engineer specializing in coding tasks, software architecture, and best practices.

**Expertise:**
- Primary Languages: Go, Python, JavaScript/TypeScript, Shell/Bash
- Frameworks: GitHub Actions, Docker, Git, REST APIs
- Specializations: Software Architecture, Code Optimization, Testing, DevOps, CI/CD

**When to use:**
- Implementing new features
- Bug fixing and debugging
- Code refactoring and optimization
- Writing and improving tests
- API development
- Script automation
- Performance improvements
- Security enhancements

**Example prompts:**
- "Implement a new feature to query ASN information from the database"
- "Optimize the database download script for better performance"
- "Add error handling to the Go query example"
- "Fix the bug in the IP address parsing logic"

**Project-Specific Context:**
- Focus: Free GeoIP databases with automatic updates
- Technologies: Go, Shell scripts, GitHub Actions, MMDB databases
- Key areas: Query tools, download automation, testing, CI/CD

### Senior Product Manager (资深产品经理)

**Files:**
- `senior-product-manager.md` (English)
- `senior-product-manager-cn.md` (中文)

**Purpose:**
The Senior Product Manager agent helps with product feature iteration and development, focusing on user-centric, practical solutions.

**Key Responsibilities:**
- Product feature iteration and planning
- User-centric approach to feature development
- Competitive analysis and inspiration from other products
- Creative problem-solving while maintaining practicality

**When to use:**
- Planning new features or improvements
- Analyzing user feedback and requirements
- Researching competitor features and best practices
- Creating product requirement documents (PRDs)
- Prioritizing feature development

**Example prompts:**
- "What features could improve the user experience for downloading databases?"
- "How do other GeoIP providers handle automatic updates?"
- "What's missing in our current documentation from a user perspective?"
- "Help me prioritize these feature requests based on user impact"

## How to Use These Agents

GitHub Copilot agents can be accessed through your IDE or GitHub interface when you have GitHub Copilot enabled. Simply reference the agent in your prompt or conversation to get specialized assistance.

Custom agents are specialized configurations that enhance GitHub Copilot's capabilities for specific tasks or domains. They provide:

1. **Domain Expertise**: Specialized knowledge in specific areas
2. **Context Awareness**: Understanding of project-specific requirements
3. **Best Practices**: Adherence to coding standards and conventions
4. **Focused Assistance**: Targeted help for particular types of tasks

## Agent Configuration Structure

For YAML-based agents (like Senior Coding Engineer), each configuration includes:

- **name**: The agent's display name
- **description**: What the agent does
- **metadata**: Version, author, and categorization information
- **expertise**: Languages, tools, and specializations
- **behavior**: Coding style, problem-solving approach, and quality standards
- **task_preferences**: Types of tasks the agent excels at
- **communication**: How the agent communicates and documents
- **project_context**: Specific knowledge about this repository
- **instructions**: Detailed guidelines for the agent's operation
- **example_prompts**: Sample prompts that work well with the agent

## Adding New Agents

To add a new agent to this project:

1. Create a new `.md` or `.yml` file in this directory
2. Define the agent's role, responsibilities, and capabilities
3. Provide clear examples of when and how to use the agent
4. Update this README to document the new agent
5. Consider creating both English and Chinese versions

## Best Practices

- **Be specific:** Clearly define the agent's domain and expertise
- **Be practical:** Focus on actionable guidance and concrete examples
- **Be clear:** Use simple language and provide clear instructions
- **Be consistent:** Follow the structure of existing agents
- Custom agents enhance but don't replace GitHub Copilot's core capabilities
- Agents work best when given clear, specific tasks
- Multiple agents can be defined for different specializations
- Keep agent configurations up to date with project changes

## Feedback

If you have suggestions for improving these agents or ideas for new agents, please:
- Open an issue in the repository
- Submit a pull request with your improvements
- Discuss in the project's discussions section
