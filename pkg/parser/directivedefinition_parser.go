package parser

import (
	"github.com/jensneuse/graphql-go-tools/pkg/document"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/token"
)

func (p *Parser) parseDirectiveDefinition(hasDescription bool, description token.Token, index *[]int) error {

	start, err := p.readExpect(keyword.DIRECTIVE, "parseDirectiveDefinition")
	if err != nil {
		return err
	}

	_, err = p.readExpect(keyword.AT, "parseDirectiveDefinition")
	if err != nil {
		return err
	}

	directiveIdent, err := p.readExpect(keyword.IDENT, "parseDirectiveDefinition")
	if err != nil {
		return err
	}

	var definition document.DirectiveDefinition
	definition.DirectiveLocations = p.indexPoolGet()
	definition.Name = p.putByteSliceReference(directiveIdent.Literal)

	if hasDescription {
		definition.Position.MergeStartIntoStart(description.TextPosition)
		definition.Description = description.Literal
	} else {
		definition.Position.MergeStartIntoStart(start.TextPosition)
	}

	err = p.parseArgumentsDefinition(&definition.ArgumentsDefinition)
	if err != nil {
		return err
	}

	_, err = p.readExpect(keyword.ON, "parseDirectiveDefinition")
	if err != nil {
		return err
	}

	for {
		next := p.l.Peek(true)

		if next == keyword.PIPE {
			p.l.Read()
		} else if next == keyword.IDENT {
			location := p.l.Read()

			parsedLocation, err := document.ParseDirectiveLocation(p.ByteSlice(location.Literal))
			if err != nil {
				return err
			}

			definition.DirectiveLocations = append(definition.DirectiveLocations, int(parsedLocation))

		} else {
			break
		}
	}

	definition.Position.MergeStartIntoEnd(p.TextPosition())
	*index = append(*index, p.putDirectiveDefinition(definition))

	return nil
}
