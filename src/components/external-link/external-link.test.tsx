import { render, screen } from "@testing-library/react";
import { describe, expect, test } from "vitest";
import { ExternalLink } from "./external-link";

describe("<ExternalLink>", () => {
  const text = "Example Link";
  const link = "https://example.com";
  test("renders without error", () => {
    render(<ExternalLink href={link}>{text}</ExternalLink>);
    expect(screen.getByText(text)).toBeInTheDocument();
  });
  test('renders an <a> tag with target="_blank" and correct rel attribute', () => {
    render(<ExternalLink href={link}>{text}</ExternalLink>);
    expect(screen.getByText(text)).toBeInstanceOf(HTMLAnchorElement);
    expect(screen.getByText(text)).toHaveAttribute("target", "_blank");
    expect(screen.getByText(text)).toHaveAttribute(
      "rel",
      "noopener noreferrer"
    );
  });
  test("accepts other attributes and attaches to underlying <a>", () => {
    render(
      <ExternalLink className="foo bar" href={link}>
        {text}
      </ExternalLink>
    );
    expect(screen.getByText(text)).toHaveClass("foo");
    expect(screen.getByText(text)).toHaveClass("bar");
  });
});
